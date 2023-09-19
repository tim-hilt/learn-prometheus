To install this chart, run the following commands from the repo-root:

``` yaml
helm dependency build ./deploy
helm install ml-simulation ./deploy # ml-simulation can be replaced with whatever
```

If you want to run all that locally, then install kind and run `kind create` before the above commands.