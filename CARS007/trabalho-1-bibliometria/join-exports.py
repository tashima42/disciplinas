from dotenv import dotenv_values

config = dotenv_values(".env")

file_count = int(config.get("WOS_FILE_COUNT"))
output_file = config.get("WOS_DATA")

with open(output_file, "w") as outfile:
    for i in range(1, file_count + 1):
        with open(f"data/savedrecs-{i}.txt", "r") as infile:
            infile_lines = infile.read().splitlines(True)
            if i > 1:
                infile_lines = infile_lines[1:] # skip the first line (header) of each file except the first one
            outfile.writelines(infile_lines)

    print(f"Joined {file_count} files into {output_file}")
