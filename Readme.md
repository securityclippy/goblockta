### WARNING ###

:warning: :warning: :warning:THIS PROJECT IS PROVIDED COMPLETELY WITHOUT WARRANTY, SUPPORT, OR ANYTHING ELSE.  USE AT YOUR OWN RISK. :warning::warning::warning:

This project is not supported by Okta in any way, shape or form




#### Summary

this is an extremely experimental project and should only be used if you know what you're doing.

It uses the Okta events API to analyze login events and looking for suspecious IP addresses which
are detected based on failure rate.


It should be noted that failure rate of an IP address can be a HORRIBLE indicator and holds the potential
to lock out a legitimate IP.

## GIANT FOOT GUN WARNING
:warning::warning: You may lock yourself out, or your entire org if you use this. :warning::warning:


THE AUTHOR OF THIS PROGRAM IS IN NO WAY RESPONSIBLE FOR ANY PROBLEMS YOU MAY CAUSE YOURSELF BY RUNNING THIS PROGRAM.  USE AT YOUR OWN RISK, you have been warned



#### getting started:


clone the repo:

```
git clone https://github.com/securityclippy/goblockta.git
sudo docker build . -t goblockta
```


edit config-example.json to match your org configuration and desired thresholds:

```json
{
  //Zones to be used for blocking
  //(Note, don't copy this, having comments in json is invalid and will cause you errors, this is just to explain teh config
  "block_zone_names": [
    "Blocked IP Addresses",
    "Blocked IP Addresses2"
  ],
  "ban_on_number_of_failures": 5,
  // time window of 5 minutes
  "ban_failure_window": 300,
  // Duration is not fully implemented
  "ban_duration": 300,
  // not implemented
  "whitelist_on_successful_mfa": true,
  "org_url": "https://myorg.okta.com",
  "okta_api_key": "long_random_okta_api_token_here",
  "log_to_slack": true,
  "slack_webhook_url": "https://hooks.slack.com/services/<long_random_string_from_slack_here>",
  "slack_user": "my_slack_user",
  //change to true to actually block.  Otherwise it will just warn you
  "blocking_mode": false,
  // warn at failure rate of 1 failure/minute
  "warn_threshold": 1.0,
  //block at filure rate of 1.5 failures/minute
  "block_threshold": 1.5
}

```

```
https://github.com/securityclippy/goblockta.git
```


