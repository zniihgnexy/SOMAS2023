import pandas as pd
def process_and_save_grouped_excel_sheets(file_path, output_file_path):
    # Read the Excel file
    xls = pd.ExcelFile(file_path)

    # Dictionary to store the results for each group across sheets
    grouped_results = {}

    # Processing each sheet
    for sheet_name in xls.sheet_names:
        # Load the sheet
        df = pd.read_excel(xls, sheet_name=sheet_name)

        # Identifying all unique group names in the column headers
        group_names = set()
        for col in df.columns:
            if "Group:" in col:
                group_name = col.split(")")[0] + ")"
                group_names.add(group_name)

        # Calculating mean for each group and updating the results
        for group in group_names:
            # Initialize group in the dictionary if not already present
            if group not in grouped_results:
                grouped_results[group] = {}

            # Filter columns for the current group
            group_columns = [col for col in df.columns if group in col]
            group_df = df[group_columns]

            # Calculate mean
            total_mean = group_df.mean().mean()

            # Update the results
            grouped_results[group][sheet_name] = total_mean

    # Converting the results to a DataFrame
    results_df = pd.DataFrame.from_dict(grouped_results, orient='index')
    results_df.reset_index(inplace=True)
    results_df.rename(columns={'index': 'Group'}, inplace=True)

    # Sorting the DataFrame by Group
    sorted_results_df = results_df.sort_values(by="Group")

    # Saving the sorted results to an Excel file
    sorted_results_df.to_excel(output_file_path, index=False)

    return output_file_path

# Example usage
input_file_path = 'statistics.xlsx'  # Your uploaded file path
output_file_path = 'grouped_sorted_means.xlsx'  # Output file path
saved_file_path = process_and_save_grouped_excel_sheets(input_file_path, output_file_path)

print(f"File saved to: {saved_file_path}")