### WARNING ###

:warning: :warning: :warning:THIS PROJECT IS PROVIDED COMPLETELY WITHOUT WARRANTY, SUPPORT, OR ANYTHING ELSE.  USE AT YOUR OWN RISK. :warning::warning::warning:




#### Summary

this is an extremely experimental project and should only be used if you know what you're doing.

It uses the Okta events API to analyze login events and looking for suspecious IP addresses which
are detected based on failure rate.


It should be noted that failure rate of an IP address can be a HORRIBLE indicator and holds the potential
to lock out a legitimate IP.



#### getting started:

from docker:
```
sudo docker build . -t goblockta
sudo docker run --rm -v $(pwd)/example-config.json:/var/config.json goblockta
```


or clone the repo:


