Change: Upgrade Scaleway SDK

As Scaleway has dropped the previously used Go library we are forced to update
to the new SDK. With this change we are forced to update the required attributes
which results in a breaking change. The new Scaleway API requires an access and
secret key instead of a single access token. Beside that the region flag had to
be dropped in favor of a zone flag.

https://github.com/promhippie/prometheus-scw-sd/pull/11
