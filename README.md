# go-timer
A implementation of hierarchical timing wheels, since build-in timer in golang has serveral limitations:

* one timer only notify once or one specific duration.
* build-in timer can't keep any states.
* build-in timer can't customize channel, **_n_** timer will create **_n_** channel.

But in many senarios, we need more from build-in timer, that's the reason **go-timer** be needed. And **go-timer** supports:

## How-to-use