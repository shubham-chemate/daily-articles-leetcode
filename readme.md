# daily-articles-leetcode

- Leetcode discuss section is gem and I don't want to miss even single good article. 
- But checking it every hour is time and mental bandwidth wasting.  
- This project fetches all the leetcode articles that are posted in last 24 hours and sends them to the email.
- Run on configured time everyday.

## How It Works

The workflow runs automatically every day at 03:45 UTC (09:15 AM IST) via GitHub Actions.

### Important Note About Scheduled Workflows

GitHub automatically disables scheduled workflows if a repository has no activity for 60 days. To prevent this:

1. **Keep-Alive Workflow**: A monthly keep-alive workflow runs on the 1st of each month to keep all workflows active.
2. **Manual Trigger**: You can manually trigger the workflow anytime from the Actions tab using the "Run workflow" button.
3. **Re-enabling**: If workflows get disabled, visit the Actions tab and re-enable them manually.

### Troubleshooting

If the workflow didn't run today:
- Check if the workflow is enabled in the Actions tab
- Manually trigger the workflow using the "Run workflow" button
- Verify that Actions are enabled in repository settings 