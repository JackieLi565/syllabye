"""
extract_programs.py

This script scrapes Toronto Metropolitan University's academic calendar to extract
program related data. The scraped data is saved as a JSON file.

Usage:
    python program_scraper.py

Arguments:
    year (optional): The academic year range to scrape (e.g., "2024-2025").
        Default: "2024-2025"

    path (optional): The output path of the JSON payload.
        Default: "./out"

    verbose (optional): Script verbose mode.
        Default: False

Outputs:
    - JSON file saved in ~/[syllabye_project]/out/programs.json containing the scraped program data.
"""

from datetime import datetime
import json
import requests
from bs4 import BeautifulSoup
import argparse

parser = argparse.ArgumentParser(description="Syllabye Program Scraper Script")

parser.add_argument(
    "--year", type=str, default="2024-2025", help="Website year range iteration"
)
parser.add_argument(
    "--path",
    type=str,
    default="./out",
    help="Directory path of output destination",
)
parser.add_argument("-v", "--verbose", action="store_true", help="Enable verbose mode")

args = parser.parse_args()


def v_print(*stm):
    """
    Helper function to print log message if verbose argument is specified.
    """
    if args.verbose:
        print(*(datetime.now().strftime("%Y-%m-%d %H:%M:%S "), *stm))


domain = "https://www.torontomu.ca"
programs_uri = f"/calendar/{args.year}/programs/"

res = requests.get(domain + programs_uri)
if res.status_code != 200:
    v_print(
        f"[request] request to {programs_uri} failed with error code {res.status_code}"
    )
    exit(f"Request Error: HTTP request to {programs_uri} failed")
else:
    v_print(f"[request] {programs_uri} request complete")

soup = BeautifulSoup(res.content, "html.parser")

# Extract program rows from the table while removing the header
table_rows = soup.find_all("tr")[1::]

programs = []

# iterate through html table and grab the fa
for row in table_rows:
    tds = row.find_all("td")
    a_tag = tds[1].find("a")

    programs.append(
        {
            "program": a_tag.text.strip(),
            "uri": "/" + "/".join(a_tag.get("href").strip().split("/")[-2:]),
            "faculty": tds[0].text.strip(),
        }
    )

output_file = f"{args.path}/courses.json"
with open(output_file, "w") as file:
    json.dump(programs, file)

v_print(f"[file] programs written to {output_file}")
