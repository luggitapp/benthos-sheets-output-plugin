# Benthos Sheets Output Plugin

This benthos plugin provides a simple way to output to a google spreadsheet.


## Getting Started
For this plugin to work you will need to get a service account credential from google cloud.

1. Go to x
2. Create the credentials files
3. Download the json key




## Plugin Configuration
For now the configuration is pretty simple.

> Paramenters
- SheetId: The sheet id of your spreesheet this can be found on the url.

> Example benthos output configuration:
```yaml
  output:
    type: sheets
    plugin:
      sheetId: "1z5W9eXXzr9fG0lYcTtXSgtYfnusIPCuAIomsZvhkO5s"
```


## Deployment
We use helm to deploy our containers with that we provide a chart for deploying your first processor to sheets.

### Usage
1. Clone this repository
2. Edit the `values.yaml` value to your needs
3. Obtain a google service account and create a secret
```kubectl create secret```
4. Install the chart `helm install myprocessor chart/ -f values.yaml`
5. Done



## Authors
> Hugo Fonseca <hugo.fonseca@luggit.app>


## LICENSE
MIT