import os
from typing import Dict, List, Tuple

import matplotlib.image as mpimg
import matplotlib.pyplot as plt
import pandas as pd
from matplotlib import cm  # colormap
from matplotlib.dates import DateFormatter, MonthLocator
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

# Convert the Month column to datetime format
data_all["Month"] = pd.to_datetime(data_all["Month"], format="%Y-%m")

# Filter out decks that have less than show_threshold 1% representation
show_threshold: float = .9
data_all = data_all[data_all["Percent"] >= show_threshold]

# Split the DataFrame into multiple DataFrames, each for a specific month
dfs_by_month: Dict[pd.Timestamp, pd.DataFrame] = {month: df for month, df in data_all.groupby("Month")}

for month, df in dfs_by_month.items():
    dfs_by_month[month] = df[df["Percent"] >= show_threshold]

# Group the data by month and sort each group by Percent
data_by_month_sorted: Dict[pd.Timestamp, pd.DataFrame] = {
    month: df.sort_values(by="Percent", ascending=False) for month, df in dfs_by_month.items()}

# Create a new figure (graph) with size in inches
plt.figure(figsize=(16, 9))
# Adjust the subplot parameters to move the plot to the left and use full vertical space
plt.subplots_adjust(left=0.05, right=0.9, top=0.95, bottom=0.05)

# Ensure all decks have entries for all months, filling missing values with 0
all_months = pd.date_range(start=data_all["Month"].min(), end=data_all["Month"].max(), freq="MS")
data_all = data_all.set_index(["Month", "Deck"]).unstack(fill_value=0).stack().reset_index()

# Calculate the average percentage for each deck
average_percentages = data_all.groupby("Deck")["Percent"].mean().sort_values(ascending=False)
print(f"average_percentages: {average_percentages}")

# Remove rows where Percent is 0
data_all = data_all[data_all["Percent"] != 0]

# seen_labels helps to avoid duplicate labels in the legend
seen_labels = set()
# previous_positions helps to adjust icon alpha if close to another icon
previous_positions: Dict[pd.Timestamp, List[Tuple[pd.Timestamp, float]]] = {}

# colormap helps to make close values have contrasting colors
colormap = cm.get_cmap('prism', len(average_percentages))

last_month = data_all["Month"].max()  # last Month in the dataset: 2024-12

# Plot each deck over the months
for idx, (month, df) in enumerate(sorted(data_by_month_sorted.items())):
    top_10_decks: List[str] = df.head(16)["Deck"].unique()
    print("_" * 16)
    print(f"month: {month.strftime('%Y-%m')}, top_10_decks: {top_10_decks}")
    if month not in previous_positions:
        previous_positions[month] = []  # Initialize list for each month
    # Find the peak percentage of all decks in the current month
    peak_of_month = df.loc[df["Percent"].idxmax()]

    for deck in top_10_decks:
        deck_data = data_all[data_all["Deck"] == deck]  # rows of a specific deck in months

        # Draw lines connecting the points
        color = colormap(average_percentages.index.get_loc(deck) / len(average_percentages))
        if deck not in seen_labels:
            line, = plt.plot(deck_data["Month"], deck_data["Percent"], label=deck, color=color)
            seen_labels.add(deck)
        else:
            line, = plt.plot(deck_data["Month"], deck_data["Percent"], color=color)
        line_color = line.get_color()  # Get the color of the line

        # Find the peak percentage of the deck over the months
        peak_of_a_deck = deck_data.loc[deck_data["Percent"].idxmax()]

        # Add icon or text at the peak of the deck over the months
        if deck in deck_icons:
            current_percent = deck_data[deck_data["Month"] == month]["Percent"].values[0]
            # only draw icon at the peak of the month, the peak of the deck, or the last month
            if current_percent in [peak_of_month["Percent"], peak_of_a_deck["Percent"]] or month == last_month:
                position = (month, current_percent)

                # if close to existed icon, adjust position and transparency
                alpha = 1.0
                for prev_pos in previous_positions[month]:
                    if position[0] == prev_pos[0] and abs(position[1] - prev_pos[1]) < 1:
                        alpha *= 0.9  # Make the image transparent if close to a previous image
                        position = (position[0], position[1] - 0.5)  # Move down
                        position = (position[0] + pd.Timedelta(days=4), position[1])  # Move to the right
                adjustedIcon = deck_icons[deck].copy()
                adjustedIcon[:, :, 3] = adjustedIcon[:, :, 3] * alpha

                imagebox = OffsetImage(adjustedIcon, zoom=0.25)
                ab = AnnotationBbox(imagebox, position, frameon=True, bboxprops=dict(
                    edgecolor='black', linewidth=0.8, boxstyle='round,pad=0'))
                plt.gca().add_artist(ab)
                previous_positions[month].append(position)
        else:
            plt.text(peak_of_a_deck["Month"], peak_of_a_deck["Percent"] + 0.5, deck, fontsize=8, ha="right")

        # Add dots for intermediate points
        for i in range(1, len(deck_data) - 1):
            row = deck_data.iloc[i]
            plt.plot(row["Month"], row["Percent"], "o", color=line_color)

# Set the labels and title
plt.xlabel("")  # "Month"
plt.ylabel("")  # "Percent"
plt.title("MasterDuelMeta 2024 Masters and DLvMax decks")
plt.grid(True)

# Set y-axis (Percent) range from 0 to 25
plt.ylim(0, 25)

# Ensure the x-axis is sorted by month
plt.gca().xaxis.set_major_formatter(DateFormatter("%Y-%m"))
plt.gca().xaxis.set_major_locator(MonthLocator())
plt.xticks(rotation=0)

# Format the y-axis to show percentage
plt.gca().yaxis.set_major_formatter(FuncFormatter(lambda y, _: f"{y}%"))

# Add legend ordered by average percentage
handles, labels = plt.gca().get_legend_handles_labels()
sorted_labels_handles = sorted(zip(labels, handles), key=lambda x: average_percentages[x[0]], reverse=True)
sorted_labels, sorted_handles = zip(*sorted_labels_handles)
plt.legend(sorted_handles, sorted_labels, title="Decks", bbox_to_anchor=(1.01, 1), loc="upper left", prop={"size": 7})

plt.show()
