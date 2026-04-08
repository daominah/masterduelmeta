import os
from typing import Dict, List, Tuple

import matplotlib
import matplotlib.image as mpimg
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from matplotlib.dates import DateFormatter, MonthLocator, date2num
from matplotlib.offsetbox import OffsetImage, AnnotationBbox
from matplotlib.ticker import FuncFormatter

from archetype_icon import archetype_icon_paths


deck_icons: Dict[str, any] = {}  # map Archetype to Image as ndarray
# Read files from archetype_icon/ using paths in archetype_icon_paths
for deck, path in archetype_icon_paths.items():
    icon_file_path = os.path.join("archetype_icon/", archetype_icon_paths[deck])
    deck_icons[deck] = mpimg.imread(icon_file_path)
    # print(f"loaded icon deck: {deck}, path: {icon_file_path}")

# Read the CSV file (created by cmd/s3_analyze_to_visualize)
data_all: pd.DataFrame = pd.read_csv("time_series_2022-01_to_2026-03.csv")

# Convert the Month column to datetime format pd.Timestamp
data_all["Month"] = pd.to_datetime(data_all["Month"], format="%Y-%m")

isLogScale = True  # isLogScale determines whether the y-axis is logarithmic or linear

top_cut_percent: float = 0.99  # filter out rows that have less than 1% representation
# draw_line_threshold: only draw a line for a deck if its peak is greater than this
# (100% meaning no lines)
draw_line_threshold = 15.00
favorite_decks = [
    # favorite decks, always draw lines for them;
    # "Bystial" is vaguely classified on MasterDuelMeta, so not included
    # "Barrier Statue",
    # "Blue-Eyes",
    # "Centur-Ion",
    # "Kashtira",
    # "Labrynth",
    # "Swordsoul",
    # "Tearlaments",
    # "True Draco",
    # "Utopia",
]

check_overlap_threshold = 1.0  # check if the icon is overlapping by y-axis difference
icon_pos_adj_y = 0.18  # adjust the y-axis position of the icon if overlapping
icon_pos_adj_x = pd.Timedelta(days=3.9)  # adjust the x-axis position of the icon if overlapping
icon_pos_adj_x2 = pd.Timedelta(days=2.4)  # adjust the x-axis but less, so icons do not overflow to the next month
zero = 0.0  # if log scale, zero needs to be a small positive number to avoid log(0) math error

# Log scale compresses the top end and stretches the bottom end.
# It is now the default for better visual distinction.
if isLogScale:
    top_cut_percent = np.log(top_cut_percent)  # convert to log space to match transformed Percent
    draw_line_threshold = np.log(draw_line_threshold)  # same
    check_overlap_threshold = 0.12
    icon_pos_adj_y = 0.02
    icon_pos_adj_x = pd.Timedelta(days=6)
    icon_pos_adj_x2 = pd.Timedelta(days=5)
    zero = 1e-4

n_months = data_all["Month"].nunique()
chart_figsize_wh = (22, 11) if n_months <= 18 else (60, 11)
dpi = 160  # dpi combined with figsize determines deck icon size relative to the whole canvas

# Ensure all decks have entries for all months so lines are continuous (no gaps).
# unstack pivots Deck into columns (filling gaps with zero), stack folds them back into rows.
all_months = pd.date_range(start=data_all["Month"].min(), end=data_all["Month"].max(), freq="MS")
data_all = data_all.set_index(["Month", "Deck"]).unstack(fill_value=zero).stack(future_stack=True).reset_index()

if isLogScale:
    data_all["Percent"] = np.log(data_all["Percent"])

data_all = data_all[data_all["Percent"] >= top_cut_percent]

# Split the DataFrame into multiple DataFrames, each for a specific month
dfs_by_month: Dict[pd.Timestamp, pd.DataFrame] = {month: df for month, df in data_all.groupby("Month")}

# Group the data by month and sort each group by Percent
data_by_month_sorted: Dict[pd.Timestamp, pd.DataFrame] = {
    month: df.sort_values(by="Percent", ascending=False) for month, df in dfs_by_month.items()}

# Create a new figure (graph) with size in inches
plt.figure(figsize=chart_figsize_wh)
# Adjust the subplot parameters to move the plot to the left and use full vertical space
plt.subplots_adjust(left=0.04, right=0.90, top=0.96, bottom=0.05)

# Calculate the average percentage for each deck
average_percentages = data_all.groupby("Deck")["Percent"].mean().sort_values(ascending=False)
print(f"average_percentages: {average_percentages.head(20)}")

# Remove fill-value rows added by unstack (zero or near-zero after log transform),
# so they do not appear as data points on the chart.
if isLogScale:
    data_all = data_all[data_all["Percent"] >= np.log(zero)]
else:
    data_all = data_all[data_all["Percent"] != 0]

# colormap helps to make close values have contrasting colors;
# sized by average_percentages so each deck gets a consistent color across months
colormap = matplotlib.cm.get_cmap('prism', len(average_percentages))

last_month: pd.Timestamp = data_all["Month"].max()  # last Month in the dataset
first_month: pd.Timestamp = data_all["Month"].min()  # first Month in the dataset

deck_peak_month = {}  # key is deck, value is the month it was used the most
decks_debut_month = {}  # key is deck, value is the first month it was in the top cut

# Calculate peak and debut months for each deck
for deck in data_all["Deck"].unique():
    deck_data = data_all[data_all["Deck"] == deck]
    deck_peak_month[deck] = deck_data.loc[deck_data["Percent"].idxmax()]
    decks_debut_month[deck] = deck_data["Month"].min()

# only draw points for the top 20 decks,
# with an exception for the last month to draw more decks
top_cut_n = 20

# previous_positions helps to adjust when icons are overlapping
previous_positions: Dict[pd.Timestamp, List[Tuple[pd.Timestamp, float]]] = {}
missing_icon_decks = set()
shown_previous_month: set = set()

# adjusted icon x per (deck, data_month); months without overlap stay at the true month
icon_adjusted_x: Dict[Tuple[str, pd.Timestamp], pd.Timestamp] = {}
# (deck, data_month, adj_x, adj_y) in draw order, so later entries render on top
icon_draw_order: List[Tuple[str, pd.Timestamp, pd.Timestamp, float]] = []
text_draw_order: List[Tuple[str, pd.Timestamp, float]] = []  # (deck, data_month, percent)
decks_ever_in_draw: set = set()  # decks that appeared in any month's decks_to_draw

# Pass 1: compute all adjusted icon positions without drawing anything
for idx, (month, df) in enumerate(sorted(data_by_month_sorted.items())):
    top_n_decks = df.head(top_cut_n)["Deck"].unique()
    if month == last_month:
        top_n_decks = df.head(top_cut_n + 4)["Deck"].unique()
    print("________________________________")
    print(f"month: {month.strftime('%Y-%m')}, top {len(top_n_decks)} decks: {top_n_decks}")
    if month not in previous_positions:
        empty_list_positions: List[Tuple[pd.Timestamp, float]] = []
        previous_positions[month] = empty_list_positions

    # Find the best decks in the current month
    tmp = df.nlargest(4, "Percent")
    month_top1 = tmp.iloc[0]
    month_top2 = tmp.iloc[1]
    month_top3 = tmp.iloc[2]
    month_top4 = tmp.iloc[3]

    decks_to_draw = top_n_decks[::-1]  # reverse order so the top deck icon is on top
    fav_decks_not_top = []  # favorite decks not in the top cut (defined in favorite_decks)
    for fav in favorite_decks:
        if fav in decks_to_draw:
            continue
        fav_row = data_all[(data_all["Deck"] == fav) & (data_all["Month"] == month)]
        if not fav_row.empty:
            percent = fav_row["Percent"].values[0]
            fav_decks_not_top.append((fav, percent))
    fav_decks_not_top.sort(key=lambda x: x[1], reverse=True)
    fav_decks_not_top = np.array([fav for fav, _ in fav_decks_not_top])
    decks_to_draw = np.concatenate((decks_to_draw, fav_decks_not_top))

    carry_decks = []  # decks shown last month, carried forward one month
    for prev_deck in shown_previous_month:
        if prev_deck in decks_to_draw:
            continue
        # even if not in top_cut_n
        prev_row = data_all[(data_all["Deck"] == prev_deck) & (data_all["Month"] == month)]
        if not prev_row.empty:
            carry_decks.append(prev_deck)
    decks_to_draw = np.concatenate((decks_to_draw, np.array(carry_decks)))

    shown_this_month: set = set()  # decks with shouldShowOrganically, used to populate shown_previous_month
    for deck in decks_to_draw:
        deck_data = data_all[data_all["Deck"] == deck]  # rows for a specific deck across months
        decks_ever_in_draw.add(deck)

        is_last_fav_appear = deck in favorite_decks and month == deck_data["Month"].max()
        current_percent = deck_data[deck_data["Month"] == month]["Percent"].values[0]
        shouldShowOrganically = (
                current_percent == month_top1["Percent"] or  # top 1 deck this month
                current_percent == month_top2["Percent"] or  # top 2 deck this month
                current_percent == month_top3["Percent"] or
                current_percent == month_top4["Percent"] or
                current_percent == deck_peak_month[deck]["Percent"] or  # current deck's all-time peak
                month == decks_debut_month[deck] or  # first month the deck appeared
                month == last_month or  # always show on the last month
                is_last_fav_appear)  # last month a favorite deck appeared
        shouldShowDeck = shouldShowOrganically or deck in shown_previous_month  # carry forward one month
        if month == first_month:  # draw fewer decks in the first month
            if (deck in top_n_decks and list(top_n_decks).index(deck) >= 10
                    and deck_peak_month[deck]["Percent"] < draw_line_threshold
                    and deck not in favorite_decks):
                shouldShowDeck = False

        if shouldShowDeck:
            if deck in deck_icons:
                # adjust position if close to an existing deck icon
                position: Tuple[pd.Timestamp, float] = (month, current_percent)
                for prev_pos in previous_positions[month]:
                    if position[0] == prev_pos[0] and abs(position[1] - prev_pos[1]) < check_overlap_threshold:
                        position = (position[0], position[1] + icon_pos_adj_y)
                        adjusted_x_axis = position[0] + icon_pos_adj_x
                        if month == last_month or month == first_month:
                            adjusted_x_axis = position[0] + icon_pos_adj_x2  # smaller shift at chart edges to avoid overflow
                        position = (adjusted_x_axis, position[1])
                icon_adjusted_x[(deck, month)] = position[0]
                icon_draw_order.append((deck, month, position[0], position[1]))
                previous_positions[month].append(position)
                if shouldShowOrganically:
                    shown_this_month.add(deck)
            else:
                missing_icon_decks.add(deck)
                text_draw_order.append((deck, month, current_percent))

    shown_previous_month = shown_this_month

print("________________________________")
for deck in missing_icon_decks:
    print(f"missing icon for deck: {deck}")

# Pass 2a: draw lines with x-endpoints adjusted to match icon centers
# Lines are drawn before icons so the icon body renders on top of the line endpoint.
seen_labels = set()
for deck in decks_ever_in_draw:
    if deck_peak_month[deck]["Percent"] >= draw_line_threshold or deck in favorite_decks:
        deck_data = data_all[data_all["Deck"] == deck]
        adjusted_months = [icon_adjusted_x.get((deck, m), m) for m in deck_data["Month"]]
        color = colormap((average_percentages.index.get_loc(deck)) / len(average_percentages))
        line_width = 0.8
        if deck not in seen_labels:
            plt.plot(adjusted_months, deck_data["Percent"], label=deck, color=color, linewidth=line_width)
            seen_labels.add(deck)
        else:
            plt.plot(adjusted_months, deck_data["Percent"], color=color, linewidth=line_width)
        # Add dots for all points except the first
        for i in range(1, len(deck_data)):
            row = deck_data.iloc[i]
            adj_month = icon_adjusted_x.get((deck, row["Month"]), row["Month"])
            plt.plot(adj_month, row["Percent"], "o", color=color, markersize=4)

# Pass 2b: draw icons on top of line endpoints
for deck, data_month, adj_x, adj_y in icon_draw_order:
    imagebox = OffsetImage(deck_icons[deck], zoom=0.25)
    drawPos = date2num(adj_x), adj_y
    ab = AnnotationBbox(imagebox, drawPos, frameon=True, bboxprops=dict(
        edgecolor='black', linewidth=0.8, boxstyle='round,pad=0'))
    plt.gca().add_artist(ab)

# Pass 2c: draw text labels for decks without icons
for deck, data_month, percent in text_draw_order:
    plt.text(data_month, percent, deck, fontsize=12)

# Set the labels and title
plt.xlabel("")  # "Month" as x-axis
plt.xlim(first_month - pd.Timedelta(days=10), last_month + pd.Timedelta(days=30))  # extend x-axis to have space for icons
plt.ylabel("")  # "Percent" as y-axis

y_axis_max_percent = 26  # for y-axis range
if isLogScale:
    plt.ylim(np.log(1.0) - 0.08, np.log(y_axis_max_percent) + 0.08)  # log
    plt.gca().yaxis.set_major_formatter(FuncFormatter(lambda y, _: f"{np.exp(y):.2f}%"))
else:
    plt.ylim(0, y_axis_max_percent)  # lin
    plt.gca().yaxis.set_major_formatter(FuncFormatter(lambda y, _: f"{y}%"))

plt.grid(True, linewidth=0.5)
plt.title(f"MasterDuelMeta {first_month.strftime('%Y-%m')} - {last_month.strftime('%Y-%m')} (Master and DLvMax decks)")

# Ensure the x-axis is sorted by month
plt.gca().xaxis.set_major_formatter(DateFormatter("%Y-%m"))
plt.gca().xaxis.set_major_locator(MonthLocator())
plt.xticks(rotation=0)

# Add legend ordered by average percentage
handles, labels = plt.gca().get_legend_handles_labels()
print(f"legend handles: n_handles: {len(handles)}, labels: {len(labels)}")
# Why is this output: legend handles: n_handles: 10, labels: 10.
# How to increase the number of handles?
# The legend only shows decks that formed a line (so drew at least twice),
# adjust this to include more decks in the legend.

if len(labels) > 0:
    sorted_labels_handles = sorted(zip(labels, handles), key=lambda x: average_percentages[x[0]], reverse=True)
    sorted_labels, sorted_handles = zip(*sorted_labels_handles)
    plt.legend(
        sorted_handles,  # plot lines (handles) to include in the legend
        sorted_labels,  # the labels corresponding to the handles
        title="Decks",  # title of the legend, this will draw: "Decks: Blue-Eyes, Tenpai Dragon, ..."
        bbox_to_anchor=(1.01, 1),  # position the legend to the plot
        loc="upper left",  # anchor the legend in the upper-left corner of the bounding box
        prop={"size": 7}  # Set the font size of the legend text to 7
    )

# Save the plot as a PNG file
date_range = f"{first_month.strftime('%Y-%m')}_{last_month.strftime('%Y-%m')}"
if isLogScale:
    plt.savefig(f"decks_log_scale_{date_range}.png", format="png", dpi=dpi)
else:
    plt.savefig(f"decks_{date_range}.png", format="png", dpi=dpi)

if n_months <= 12:
    plt.show()
