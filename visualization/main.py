import os
import matplotlib.image as mpimg
import matplotlib.pyplot as plt
import pandas as pd

from archetype_icon import archetype_icon_paths

# Read the CSV file
df = pd.read_csv('time_series_2024.csv')

# Filter to keep only the top 10 decks for each month
df_top10 = df.groupby('Month').apply(lambda x: x.nlargest(10, 'Percent')).reset_index(drop=True)

# Pivot the data to have months as rows and decks as columns
df_pivot = df_top10.pivot(index='Month', columns='Deck', values='Percent')

# Plot the data
plt.figure(figsize=(10, 6))
for column in df_pivot.columns:
    plt.plot(df_pivot.index, df_pivot[column], label=column)

# Set the labels and title
plt.xlabel('Month')
plt.ylabel('Percent')
plt.title('Top 10 Archetype Distribution Over Time')
plt.legend(loc='best')
plt.grid(True)

# Add icons to the legend
handles, labels = plt.gca().get_legend_handles_labels()
new_handles = []
for label in labels:
    if label in archetype_icon_paths:
        # Load the image, icon_path is in dictionary archetype_icon_paths
        icon_path = os.path.join("archetype_icon/", archetype_icon_paths[label])
        img = mpimg.imread(icon_path)
        new_handles.append(plt.Line2D([0], [0], marker='o', color='w', markerfacecolor='w', markersize=10,
                                      markeredgecolor='w', markeredgewidth=0, label=label))
    else:
        new_handles.append(handles[labels.index(label)])

plt.legend(handles=new_handles, loc='best')

plt.show()
