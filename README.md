# pipeline-queue

A naive solution to ensuring a singleton pipeline queue for GitLab pipelines.

## Usage

```
# full command with default values as explicit options
$ pipeline-queue --hostname https://gitlab.com --interval-time 30s -pipeline $CI_PIPELINE_ID --project $CI_PROJECT_ID --token $GITLAB_API_TOKEN

# full command with using shorthand flags
$ pipeline-queue -n https://gitlab.com -i 30s -l $CI_PIPELINE_ID -j $CI_PROJECT_ID -t $GITLAB_API_TOKEN

# equivalent abbreviated version of the above
$ pipeline-queue -t $GITLAB_API_TOKEN
```