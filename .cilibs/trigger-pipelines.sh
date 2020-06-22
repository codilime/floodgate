#!/bin/bash -e
for PIPELINE in $@ ; do

  PASS=`cat ~/.hal/default/profiles/gate-local.yml | grep password`
  PASS=${PASS#*:\ }
  USER=`cat ~/.hal/default/profiles/gate-local.yml | grep name`
  USER=${USER#*:\ }
  ALL_APPS=`curl -s -X GET --user "$USER:$PASS" "http://spinnaker/api/v1/applications" | jq -r .[].name`
  MAX_ATTEMPTS=20
  
  echo "Triggering pipeline with source $PIPELINE"
  EVENT_ID=`curl -s -X POST -H "content-type: application/json" -d "{ }" http://spinnaker/api/v1/webhooks/webhook/$PIPELINE | jq -r .eventId`
  echo "eventId: $EVENT_ID"
  
  for APP in $ALL_APPS ; do
    PIPELINE_NAME=`curl -s -X GET --user "$USER:$PASS" "http://spinnaker/api/v1/applications/$APP/executions/search?triggerTypes=webhook&eventId=$EVENT_ID" | jq -r .[].name`
    ATTEMPTS=0

    while [[ $PIPELINE_NAME != "" ]] &&  [ $ATTEMPTS -lt $MAX_ATTEMPTS ] ; do
      echo "Checking pipeline $PIPELINE_NAME status"
      STATUS=`curl -s -X GET --user "$USER:$PASS" "http://spinnaker/api/v1/applications/$APP/executions/search?triggerTypes=webhook&eventId=$EVENT_ID" | jq -r .[].status`

      case $STATUS in

        "NOT_STARTED")
          echo "Waiting for pipeline $PIPELINE_NAME to start"
          sleep 3
	  ;;

        "RUNNING")
          echo "Waiting for pipeline $PIPELINE_NAME to finish"
          sleep 3
	  ;;

        "SUCCEEDED")
          echo "$Pipeline PIPELINE_NAME succeded"
          break
          ;;

        *)
          echo "Pipeline $PIPELINE_NAME exited with status $STATUS"
          exit 1
	  ;;
      esac
      ((++ATTEMPTS))
    done

    if [ $ATTEMPTS -ge $MAX_ATTEMPTS ] ; then
      echo "Check timed out"
      exit 1
    fi
  done
done
