#!/bin/bash

WINDOW_SIZE=1 # constant load window size in seconds
CPU_LOAD=1 # how many cores to use
DIST="conts" # distributon. one of conts|random|uniform|[10,20,0,50,5,15,0]

for i in "$@"; do
case ${i} in
    --window-size=*)
        WINDOW_SIZE=${i#*=}
        shift # past argument=value
        ;;

    --cpu-load=*)
        CPU_LOAD=${i#*=}
        shift
        ;;

    --dist=*)
        ENV_DIR="${i#*=}"
        ENV_DIR="${ENV_DIR/#\~/$HOME}"
        if [[ ! -d ${ENV_DIR} ]]; then
            echo "${ENV_DIR} is not directory"
            exit 1
        fi

        TMP_FILE=${ENV_DIR}/.aws.env
        cat ${ENV_DIR}/credentials > ${TMP_FILE}
        cat ${ENV_DIR}/config >> ${TMP_FILE}

        sed -i 's/\[default\]//g' ${TMP_FILE}
        sed -i 's/ = /=/g' ${TMP_FILE}
        sed -i 's/aws_access_key_id/AWS_ACCESS_KEY_ID/g' ${TMP_FILE}
        sed -i 's/aws_secret_access_key/AWS_SECRET_ACCESS_KEY/g' ${TMP_FILE}
        sed -i 's/region/AWS_DEFAULT_REGION/g' ${TMP_FILE}

        RUN_FLAGS="${RUN_FLAGS} --env-file ${TMP_FILE}"
        CLD_PROVIDER=1

        shift # past argument=value
        ;;

    -g|--gcp)
        CLD_PROVIDER=2
        shift # past argument
        ;;

    *)
        echo "Invalid usage"
        exit 1
esac
done
