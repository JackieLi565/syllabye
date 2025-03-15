"""
extract_courses.py

This script scrapes Toronto Metropolitan University's academic calendar to extract
course related data. The scraped data is saved as a JSON file.

Usage:
    python load_data.py

Arguments:
    env-file (optional): Environment variables file.
        Default: ".env"

    verbose (optional): Enable verbose mode.
        Default: False

    path (optional): Directory path of JSON source files.
        Default: "./out"
"""

import json
import psycopg2 as pg
from psycopg2.extras import execute_values
from datetime import datetime
import argparse
from dotenv import dotenv_values

parser = argparse.ArgumentParser(description="Syllabye Data Loader Script")

parser.add_argument("--env-file", type=str, default=".env", help="Postgres host")

parser.add_argument("-v", "--verbose", action="store_true", help="Enable verbose mode")
parser.add_argument(
    "--path",
    type=str,
    default="./out",
    help="Directory path of json source files",
)

args = parser.parse_args()


def v_print(*stm):
    """Helper function to print log message if verbose argument is specified."""
    if args.verbose:
        print(*(datetime.now().strftime("%Y-%m-%d %H:%M:%S "), *stm))


try:
    env_vars = dotenv_values(args.env_file)

    conn = pg.connect(
        host=env_vars["POSTGRES_HOST"],
        port=env_vars["POSTGRES_PORT"],
        dbname=env_vars["POSTGRES_DATABASE"],
        user=env_vars["POSTGRES_USER"],
        password=env_vars["POSTGRES_PASSWORD"],
    )
    cur = conn.cursor()
    v_print("[database] Connection established")
except Exception as e:
    exit(e)
    v_print(
        (
            "[database] connection failed with credential values:\n\n"
            f"\thost={args.host}\n"
            f"\tport={args.port}\n"
            f"\tdatabase={args.database}\n"
            f"\tuser={args.user}\n"
            f"\tpassword={args.password}\n"
        )
    )

programs_path = f"{args.path}/programs.json"
with open(programs_path) as program_file:
    programs = json.load(program_file)
    v_print(f"[file] loaded {programs_path}")

faculties = set([program["faculty"] for program in programs])

execute_values(
    cur=cur,
    sql=r"""
insert into faculties (name)
values %s
returning id, name;
""",
    argslist=[(faculty,) for faculty in faculties],
)

inserted_faculties = cur.execute(
    r"""
select id, name
from faculties;
"""
)
faculty_map = {name: faculty_id for faculty_id, name in cur.fetchall()}

# Insert programs
execute_values(
    cur=cur,
    sql=r"""
insert into programs (faculty_id, name, uri)
values %s
""",
    argslist=[
        (faculty_map[program["faculty"]], program["program"], program["uri"])
        for program in programs
    ],
)

courses_path = f"{args.path}/courses.json"
with open(courses_path) as course_file:
    courses = json.load(course_file)
    v_print(f"[file] loaded {courses_path}")

course_categories = set([course["category"] for course in courses])

# Insert course_categories
execute_values(
    cur=cur,
    sql=r"""
insert into course_categories (name)
values %s;
""",
    argslist=[(category,) for category in course_categories],
)

inserted_course_categories = cur.execute(
    r"""
select id, name
from course_categories;
"""
)
course_category_map = {name: category_id for category_id, name in cur.fetchall()}

# Insert courses
execute_values(
    cur=cur,
    sql=r"""
insert into courses (category_id, title, description, uri, course, alpha, code)
values %s;
""",
    argslist=[
        (
            course_category_map[course["category"]],
            course["title"],
            course["description"],
            course["uri"],
            course["course"],
            course["alpha"],
            course["code"],
        )
        for course in courses
    ],
)

conn.commit()
v_print("[database] transaction committed with no issues")

cur.close()
conn.close()
v_print("[database] connection closed")
