# impenetrable
Basically a one-liner for creating cluster-wide sealed secrets

Assumes you have the `kubeseal` client installed.

Assumes you have the sealed secrets controller on your cluster.

It doesn't care about your cert. `kubeseal` pulls your the controller's cert automatically.

Cluster-wide because it's the easy route.

Just enter plain text and get back the encrypted secret.

```
impenetrable is-your-sweet-sweet-secret
```

## Installing

Download the release binary, then extract it. e.g.:

```
tar -xzvf Linux_x86_64.tar.gz
```

Move it to your favorite local bin, e.g.:

```
sudo mv ./Linux_x86_64 /usr/local/bin/impenetrable
```
