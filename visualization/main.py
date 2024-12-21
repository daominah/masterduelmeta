import os

import matplotlib.image as mpimg
import matplotlib.pyplot as plt
import pandas as pd
from matplotlib.offsetbox import OffsetImage, AnnotationBbox

from archetype_icon import archetype_icon_paths

deck_icons = {}
# read file in dir archetype_icon/, path archetype_icon_paths[deck]
for deck, path in archetype_icon_paths.items():
    icon_file_path = os.path.join("archetype_icon/", archetype_icon_paths[deck])
    deck_icons[deck] = mpimg.imread(icon_file_path)
    print(f"loaded icon deck: {deck}, path: {icon_file_path}")

# Read the CSV file
df_all = pd.read_csv("time_series_2024.csv")

# Convert the Month column to datetime format
df_all['Month'] = pd.to_datetime(df_all['Month'], format='%Y-%m')

# Filter out decks that have less than show_threshold percent
show_threshold = 1
df_all = df_all[df_all['Percent'] >= show_threshold]

# Split the DataFrame into multiple DataFrames, each for a specific month
dfs_by_month = {month: df for month, df in df_all.groupby("Month")}

for month, df in dfs_by_month.items():
    dfs_by_month[month] = df[df['Percent'] >= show_threshold]

# Group the data by month and sort each group by Percent
dfs_by_month_sorted = {month: df.sort_values(by='Percent', ascending=False) for month, df in dfs_by_month.items()}

# Create a new figure (graph) with size in inches
plt.figure(figsize=(16, 9))

# Plot each deck over the months
for month, df in sorted(dfs_by_month_sorted.items()):
    top_10_decks = df.head(10)['Deck'].unique()
    for deck in top_10_decks:
        deck_data = df_all[df_all['Deck'] == deck]
        line, = plt.plot(deck_data['Month'], deck_data['Percent'], label=deck)  # Draw lines connecting the points
        line_color = line.get_color()  # Get the color of the line

        # Add icon or text at the beginning of the line
        first_row = deck_data.iloc[0]
        if deck in deck_icons:
            imagebox = OffsetImage(deck_icons[deck], zoom=0.16)  # Adjust zoom icon size suitably
            ab = AnnotationBbox(imagebox, (first_row['Month'], first_row['Percent'] + 0.5), frameon=False)  # Adjust position to avoid overlap
            plt.gca().add_artist(ab)
        else:
            plt.text(first_row['Month'], first_row['Percent'] + 0.5, deck, fontsize=8, ha='right')

        # Add icon or text at the end of the line
        last_row = deck_data.iloc[-1]
        if deck in deck_icons:
            imagebox = OffsetImage(deck_icons[deck], zoom=0.16)  # Adjust zoom icon size suitably
            ab = AnnotationBbox(imagebox, (last_row['Month'], last_row['Percent'] + 0.5), frameon=False)  # Adjust position to avoid overlap
            plt.gca().add_artist(ab)
        else:
            plt.text(last_row['Month'], last_row['Percent'] + 0.5, deck, fontsize=8, ha='right')

        # Add dots for intermediate points
        for i in range(1, len(deck_data) - 1):
            row = deck_data.iloc[i]
            plt.plot(row['Month'], row['Percent'], 'o', color=line_color)

# Set the labels and title
plt.xlabel('Month')
plt.ylabel('Percent')
plt.title('Master Duel meta 2024')
plt.grid(True)

# Set x-axis (Percent) range from 0 to 25
plt.ylim(0, 25)

# Ensure the x-axis is sorted by month
plt.gca().xaxis.set_major_formatter(plt.matplotlib.dates.DateFormatter('%Y-%m'))
plt.gca().xaxis.set_major_locator(plt.matplotlib.dates.MonthLocator())
plt.xticks(rotation=45)
plt.show()
