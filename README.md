# arssh

SSH into AWS instances by name tag or instance id.

Features:

- Caches the instance list so subsequent sshes don't take too long.
    - There's a random based refresh so it doesn't go stale forever.
    - You can force a refresh by passing `CACHEBUST=1` as an env var. I wanted to keep the command line as close to ssh as possible so I decided on an env var.
- Tries to guess the default user based on the AMI name (only guesses `ubuntu` and the rest uses `ec2-user`)
- Any extra arguments you would normally pass to ssh are also sent.