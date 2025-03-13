"""
extract_programs.py

This script scrapes Toronto Metropolitan University's academic calendar to extract
program related data. The scraped data is saved as a JSON file.

Usage:
    python program_scraper.py [year_range]

Arguments:
    year_range (optional): The academic year range to scrape (e.g., "2024-2025").
                           Defaults to "2024-2025" if not provided.

Outputs:
    - JSON file saved in ./out/[year_range]_programs.json containing the scraped program data.

Example:
    python program_scraper.py 2023-2024
"""

import sys
import json
import requests
from bs4 import BeautifulSoup

arguments = sys.argv

year_range = "2024-2025"
if len(arguments) == 2:
    year_range = arguments[1]

domain = "https://www.torontomu.ca"
programs_uri = f"/calendar/{year_range}/programs/"

res = requests.get(domain + programs_uri)
if res.status_code != 200:
    exit(f"Request error code: {res.status_code}")

soup = BeautifulSoup(res.content, "html.parser")

# Extract program rows from the table while removing the header
table_rows = soup.find_all("tr")[1::]

programs = []

# iterate through html table and grab the fa
for row in table_rows:
    tds = row.find_all("td")
    a_tag = tds[1].find("a")

    faculty = tds[0].text.strip()
    program_name = a_tag.text.strip()
    href = a_tag.get("href").strip()

    programs.append(
        {"program": program_name, "href": domain + href, "faculty": faculty}
    )

output_file = f"./out/{year_range}_programs.json"
with open(output_file, "w") as file:
    json.dump(programs, file)

print(f"Programs successfully saved to {output_file}")
