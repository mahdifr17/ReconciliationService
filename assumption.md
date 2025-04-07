# Assumption:
- Transaction data, both internal and bank statements are ordered by transaction date
- There is no definition of user interface on the problem statement, so this service will utilize a CLI
- Related start-end time
  - start time is 00:00:00
  - end time is 59:59:59
  - Since bank statement transaction date only yyyy-mm-dd. It will be treated as minute 5 to handle start-end filter