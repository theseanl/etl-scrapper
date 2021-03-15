# ETL Course Material crawler

Crawls course materials from ETL course websites.

## Usage

```
snuetl --username <mysnu_username> --password <mysnu_password> --course <etl_course_id>
```
Then it will download course materials to the current working directory. It will check if a file with a name it is about to create already exists and skips, so one can resume downloading by running it again.
