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
# read file in dir archetype_icon/, path archetype_icon_paths[deck]
for deck, path in archetype_icon_paths.items():
    icon_file_path = os.path.join("archetype_icon/", archetype_icon_paths[deck])
    deck_icons[deck] = mpimg.imread(icon_file_path)
    # print(f"loaded icon deck: {deck}, path: {icon_file_path}")

# Read the CSV file
data_all: pd.DataFrame = pd.read_csv("time_series_2024.csv")

# Convert the Month column to datetime format pd.Timestamp
data_all["Month"] = pd.to_datetime(data_all["Month"], format="%Y-%m")

isLogScale = True  # isLogScale determines whether the y-axis is logarithmic or linear

top_cut_percent: float = 1.00  # filter out rows that have less than 1% representation
draw_line_threshold = 10.00  # only draw lines for a deck if its peak is greater than this value
favorite_decks = [  # favorite decks, always draw lines for them; "Bystial" is vaguely classified on MasterDuelMeta, so not included
    "Barrier Statue", "Kashtira", "Centur-Ion", "Labrynth", "Tearlaments",
    "True Draco", "Swordsoul", "Utopia", "Blue-Eyes", ]
check_overlap_threshold = 1.0  # check if the icon is overlapping by y-axis difference
icon_pos_adj_y = 0.18  # adjust the y-axis position of the icon if overlapping
icon_pos_adj_x = pd.Timedelta(days=3.9)  # adjust the x-axis position of the icon if overlapping
icon_pos_adj_x2 = pd.Timedelta(days=2.4)  # adjust the x-axis but less, so many icons do not go to the next month
zero = 0.0  # if log scale, zero needs to be a small positive number to avoid log(0) math error

if isLogScale:
    top_cut_percent = np.log(1.00)
    draw_line_threshold = np.log(10.00)
    check_overlap_threshold = 0.1
    icon_pos_adj_y = 0.02
    icon_pos_adj_x = pd.Timedelta(days=5)
    icon_pos_adj_x2 = pd.Timedelta(days=5)
    zero = 1e-4

# Ensure all decks have entries for all months, filling missing values with 0
all_months = pd.date_range(start=data_all["Month"].min(), end=data_all["Month"].max(), freq="MS")
data_all = data_all.set_index(["Month", "Deck"]).unstack(fill_value=zero).stack().reset_index()

if isLogScale:
    data_all["Percent"] = np.log(data_all["Percent"])

data_all = data_all[data_all["Percent"] >= top_cut_percent]

# Split the DataFrame into multiple DataFrames, each for a specific month
dfs_by_month: Dict[pd.Timestamp, pd.DataFrame] = {month: df for month, df in data_all.groupby("Month")}

# Group the data by month and sort each group by Percent
data_by_month_sorted: Dict[pd.Timestamp, pd.DataFrame] = {
    month: df.sort_values(by="Percent", ascending=False) for month, df in dfs_by_month.items()}

# Create a new figure (graph) with size in inches
plt.figure(figsize=(18, 12))
# Adjust the subplot parameters to move the plot to the left and use full vertical space
plt.subplots_adjust(left=0.04, right=0.90, top=0.96, bottom=0.05)

# Calculate the average percentage for each deck
average_percentages = data_all.groupby("Deck")["Percent"].mean().sort_values(ascending=False)
print(f"average_percentages: {average_percentages.head(20)}")

# Remove rows where Percent is 0
if isLogScale:
    data_all = data_all[data_all["Percent"] >= np.log(zero)]
else:
    data_all = data_all[data_all["Percent"] != 0]

# seen_labels helps to avoid duplicate labels in the legend
seen_labels = set()
# previous_positions helps to adjust when icons are overlapping
previous_positions: Dict[pd.Timestamp, List[Tuple[pd.Timestamp, float]]] = {}

# colormap helps to make close values have contrasting colors
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
# can have exception for the last month to draw more decks
top_cut_n = 20

missing_icon_decks = set()
# Plot each deck over the months
for idx, (month, df) in enumerate(sorted(data_by_month_sorted.items())):
    top_n_decks = df.head(top_cut_n)["Deck"].unique()
    if month == last_month:
        top_n_decks = df.head(top_cut_n + 2)["Deck"].unique()
    print("________________________________")
    print(f"month: {month.strftime('%Y-%m')}, top {len(top_n_decks)} decks: {top_n_decks}")
    if month not in previous_positions:
        empty_list_positions: List[Tuple[pd.Timestamp, float]] = []
        previous_positions[month] = empty_list_positions

    # Find the best decks in the current month
    tmp = df.nlargest(2, "Percent")
    month_top1 = tmp.iloc[0]
    month_top2 = tmp.iloc[1]

    decks_to_draw = top_n_decks[::-1]  # reverse order so the top deck icon is on top
    fav_decks_not_top = []  # draw for favorite decks not in top cut
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

    for deck in decks_to_draw:
        deck_data = data_all[data_all["Deck"] == deck]  # rows of a specific deck in months
        # Draw lines connecting the points for notable decks
        if (deck_peak_month[deck]["Percent"] >= draw_line_threshold
                or deck in favorite_decks):
            color = colormap(average_percentages.index.get_loc(deck) / len(average_percentages))
            line_width = 0.8
            if deck not in seen_labels:
                line, = plt.plot(deck_data["Month"], deck_data["Percent"], label=deck, color=color, linewidth=line_width)
                seen_labels.add(deck)
            else:
                line, = plt.plot(deck_data["Month"], deck_data["Percent"], color=color, linewidth=line_width)

            # Add dots for intermediate points
            for i in range(1, len(deck_data) - 1):
                row = deck_data.iloc[i]
                plt.plot(row["Month"], row["Percent"], "o", color=line.get_color(), markersize=4)

        is_last_fav_appear = deck in favorite_decks and month == deck_data["Month"].max()
        current_percent = deck_data[deck_data["Month"] == month]["Percent"].values[0]
        shouldShowDeck = (current_percent == month_top1["Percent"] or
                          current_percent == month_top2["Percent"] or
                          current_percent == deck_peak_month[deck]["Percent"] or
                          month == decks_debut_month[deck] or
                          month == last_month or
                          is_last_fav_appear)
        if month == first_month:  # draw less decks in the first month
            if (list(top_n_decks).index(deck) >= 10
                    and deck_peak_month[deck]["Percent"] < draw_line_threshold
                    and deck not in favorite_decks):
                shouldShowDeck = False

        if shouldShowDeck:
            if deck in deck_icons:
                # adjust position if close to existed deck icon
                position: Tuple[pd.Timestamp, float] = (month, current_percent)
                for prev_pos in previous_positions[month]:
                    if position[0] == prev_pos[0] and abs(position[1] - prev_pos[1]) < check_overlap_threshold:
                        position = (position[0], position[1] + icon_pos_adj_y)
                        adjusted_x_axis = position[0] + icon_pos_adj_x
                        if month == last_month or month == first_month:
                            adjusted_x_axis = position[0] + icon_pos_adj_x2
                        position = (adjusted_x_axis, position[1])

                imagebox = OffsetImage(deck_icons[deck], zoom=0.25)
                drawPos = date2num(position[0]), position[1]
                ab = AnnotationBbox(imagebox, drawPos, frameon=True, bboxprops=dict(
                    edgecolor='black', linewidth=0.8, boxstyle='round,pad=0'))
                plt.gca().add_artist(ab)
                previous_positions[month].append(position)
            else:
                missing_icon_decks.add(deck)
                plt.text(month, current_percent, deck, fontsize=12)

print("________________________________")
for deck in missing_icon_decks:
    print(f"missing icon for deck: {deck}")

# Set the labels and title
plt.xlabel("")  # "Month" as x-axis
plt.xlim(first_month - pd.Timedelta(days=10), last_month + pd.Timedelta(days=30))  # extend x-axis to have space for icons
plt.ylabel("")  # "Percent" as y-axis

y_axis_max_percent = 21.5  # for y-axis range
if isLogScale:
    plt.ylim(np.log(1.0) - 0.08, np.log(y_axis_max_percent) + 0.08)  # log
    plt.gca().yaxis.set_major_formatter(FuncFormatter(lambda y, _: f"{np.exp(y):.2f}%"))
else:
    plt.ylim(0, y_axis_max_percent)  # lin
    plt.gca().yaxis.set_major_formatter(FuncFormatter(lambda y, _: f"{y}%"))

plt.grid(True, linewidth=0.5)
plt.title("MasterDuelMeta 2024 (Master and DLvMax decks)")

# Ensure the x-axis is sorted by month
plt.gca().xaxis.set_major_formatter(DateFormatter("%Y-%m"))
plt.gca().xaxis.set_major_locator(MonthLocator())
plt.xticks(rotation=0)

# Add legend ordered by average percentage
handles, labels = plt.gca().get_legend_handles_labels()
if len(labels) > 0:
    sorted_labels_handles = sorted(zip(labels, handles), key=lambda x: average_percentages[x[0]], reverse=True)
    sorted_labels, sorted_handles = zip(*sorted_labels_handles)
    plt.legend(sorted_handles, sorted_labels, title="Decks",
               bbox_to_anchor=(1.02, 1), loc="upper left", prop={"size": 7})

# Save the plot as a PNG file
if isLogScale:
    plt.savefig("decks_2024_log_scale.png", format="png", dpi=200)
else:
    plt.savefig("decks_2024.png", format="png", dpi=200)

plt.show()
