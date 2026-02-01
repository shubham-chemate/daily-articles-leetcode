# daily-articles-leetcode

- Leetcode discuss section is gem and I don't want to miss even single good article. 
- But checking it every hour is time and mental bandwidth wasting.  
- This project fetches all the leetcode articles that are posted in last 24 hours and sends them to the email.
- Run on configured time everyday.

## How It Works

The workflow runs automatically every day at 03:45 UTC (09:15 AM IST) via GitHub Actions.

Since the workflow runs daily, it naturally prevents GitHub's 60-day inactivity timeout that would otherwise disable scheduled workflows.

### Troubleshooting

If the workflow didn't run today:
- Check if the workflow is enabled in the Actions tab (it may have been disabled due to previous inactivity)
- Manually trigger the workflow using the "Run workflow" button
- Verify that Actions are enabled in repository settings
- Once running, the daily schedule will keep the workflow active 