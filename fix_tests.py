
import re
import sys

files = [
    "internal/application/coin_service_test.go",
    "internal/application/coin_service_coverage_test.go",
    "internal/application/coin_service_more_coverage_test.go"
]

for filepath in files:
    with open(filepath, 'r') as f:
        lines = f.readlines()
    
    new_lines = []
    for line in lines:
        if ":= setupTest(t)" in line:
            parts = line.split(":=")
            lhs = parts[0]
            commas = lhs.count(',')
            # We want 9 items, so 8 commas.
            # If we have 7 commas (8 items), add one.
            if commas == 7:
                 # Check if line ends with space
                 prefix = lhs.rstrip()
                 new_line = prefix + ", _ :=" + parts[1]
                 new_lines.append(new_line)
                 print(f"Fixed {filepath}: {line.strip()} -> {new_line.strip()}")
            elif commas == 8:
                 new_lines.append(line)
            else:
                 print(f"Skipping {filepath}: {line.strip()} (commas={commas})")
                 new_lines.append(line)
        else:
            new_lines.append(line)
            
    with open(filepath, 'w') as f:
        f.writelines(new_lines)
