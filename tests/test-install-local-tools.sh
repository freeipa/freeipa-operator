#!/bin/bash

DEFAULT_TIMEOUT=1200   # 20 minutes

declare -a MATRIX_FLAGS

MATRIX_FLAGS[1]="NO NO NO"
MATRIX_FLAGS[2]="YES NO NO"
MATRIX_FLAGS[3]="ASK NO NO"
MATRIX_FLAGS[4]="ASK NO NO"

MATRIX_FLAGS[5]="NO YES NO"
MATRIX_FLAGS[6]="YES YES NO" #
MATRIX_FLAGS[7]="ASK YES NO"
MATRIX_FLAGS[8]="ASK YES NO" #

MATRIX_FLAGS[9]="NO ASK NO"
MATRIX_FLAGS[10]="YES ASK NO"
MATRIX_FLAGS[11]="ASK ASK NO"
MATRIX_FLAGS[12]="ASK ASK NO"
MATRIX_FLAGS[13]="NO ASK NO"
MATRIX_FLAGS[14]="YES ASK NO"
MATRIX_FLAGS[15]="ASK ASK NO" #
MATRIX_FLAGS[16]="ASK ASK NO" #

MATRIX_FLAGS[17]="NO ASK YES"
MATRIX_FLAGS[18]="YES ASK YES"
MATRIX_FLAGS[19]="ASK ASK YES"
MATRIX_FLAGS[20]="ASK ASK YES"
MATRIX_FLAGS[21]="NO ASK YES"   #
MATRIX_FLAGS[22]="YES ASK YES"  #
MATRIX_FLAGS[23]="ASK ASK YES"  #
MATRIX_FLAGS[24]="ASK ASK YES"  #

MATRIX_FLAGS[25]="NO ASK ASK"
MATRIX_FLAGS[26]="YES ASK ASK"
MATRIX_FLAGS[27]="ASK ASK ASK"
MATRIX_FLAGS[28]="ASK ASK ASK"
MATRIX_FLAGS[29]="NO ASK ASK"
MATRIX_FLAGS[30]="YES ASK ASK"
MATRIX_FLAGS[31]="ASK ASK ASK"
MATRIX_FLAGS[32]="ASK ASK ASK"
MATRIX_FLAGS[33]="NO ASK ASK"
MATRIX_FLAGS[34]="YES ASK ASK"
MATRIX_FLAGS[35]="ASK ASK ASK"
MATRIX_FLAGS[36]="ASK ASK ASK"
MATRIX_FLAGS[37]="NO ASK ASK"
MATRIX_FLAGS[38]="YES ASK ASK"
MATRIX_FLAGS[39]="ASK ASK ASK"
MATRIX_FLAGS[40]="ASK ASK ASK"

MATRIX_FLAGS[41]="YES NO NO"


declare -a MATRIX_INPUTS

MATRIX_INPUTS[1]=""
MATRIX_INPUTS[2]=""
MATRIX_INPUTS[3]="n"
MATRIX_INPUTS[4]="y"

MATRIX_INPUTS[5]=""
MATRIX_INPUTS[6]="" #
MATRIX_INPUTS[7]="n"
MATRIX_INPUTS[8]="y" #

MATRIX_INPUTS[9]="n"
MATRIX_INPUTS[10]="n"
MATRIX_INPUTS[11]="n"
MATRIX_INPUTS[12]="yn"
MATRIX_INPUTS[13]=""
MATRIX_INPUTS[14]="n"
MATRIX_INPUTS[15]="n"  #
MATRIX_INPUTS[16]="yn" #

MATRIX_INPUTS[17]=""
MATRIX_INPUTS[18]="n"
MATRIX_INPUTS[19]="n"
MATRIX_INPUTS[20]="yn"
MATRIX_INPUTS[21]=""    #
MATRIX_INPUTS[22]="y"   #
MATRIX_INPUTS[23]="n"   #
MATRIX_INPUTS[24]="yn"  #

MATRIX_INPUTS[25]="n"
MATRIX_INPUTS[26]="n"
MATRIX_INPUTS[27]="nn"
MATRIX_INPUTS[28]="ynn"
MATRIX_INPUTS[29]="n"
MATRIX_INPUTS[30]="nn"
MATRIX_INPUTS[31]="nn"
MATRIX_INPUTS[32]="yny"
MATRIX_INPUTS[33]="y"
MATRIX_INPUTS[34]="ny"
MATRIX_INPUTS[35]="ny"
MATRIX_INPUTS[36]="yny"
MATRIX_INPUTS[37]="y"
MATRIX_INPUTS[38]="ny"
MATRIX_INPUTS[39]="ny"
MATRIX_INPUTS[40]="yny"

MATRIX_INPUTS[41]=""


declare -a MATRIX_EXPECTED

MATRIX_EXPECTED[1]=0
MATRIX_EXPECTED[2]=0
MATRIX_EXPECTED[3]=0
MATRIX_EXPECTED[4]=0

MATRIX_EXPECTED[5]=0
MATRIX_EXPECTED[6]=0
MATRIX_EXPECTED[7]=0
MATRIX_EXPECTED[8]=0

MATRIX_EXPECTED[9]=0
MATRIX_EXPECTED[10]=0
MATRIX_EXPECTED[11]=0
MATRIX_EXPECTED[12]=0
MATRIX_EXPECTED[13]=0
MATRIX_EXPECTED[14]=0
MATRIX_EXPECTED[15]=0
MATRIX_EXPECTED[16]=0

MATRIX_EXPECTED[17]=0
MATRIX_EXPECTED[18]=0
MATRIX_EXPECTED[19]=0
MATRIX_EXPECTED[20]=0
MATRIX_EXPECTED[21]=0
MATRIX_EXPECTED[22]=0
MATRIX_EXPECTED[23]=0
MATRIX_EXPECTED[24]=0

MATRIX_EXPECTED[25]=0
MATRIX_EXPECTED[26]=0
MATRIX_EXPECTED[27]=0
MATRIX_EXPECTED[28]=0
MATRIX_EXPECTED[29]=0
MATRIX_EXPECTED[30]=0
MATRIX_EXPECTED[31]=0
MATRIX_EXPECTED[32]=0
MATRIX_EXPECTED[32]=0
MATRIX_EXPECTED[34]=0
MATRIX_EXPECTED[35]=0
MATRIX_EXPECTED[36]=0
MATRIX_EXPECTED[37]=0
MATRIX_EXPECTED[38]=0
MATRIX_EXPECTED[39]=0
MATRIX_EXPECTED[40]=0

MATRIX_EXPECTED[41]=0


declare -a CONTAINER_IMAGES

# shellcheck disable=SC2153
if [ "${CONTAINER_IMAGE}" == "" ]
then
    CONTAINER_IMAGES=("centos:8" \
                      "fedora:30" \
                      "fedora:31" \
                      "fedora:32" \
                      "debian:10" \
                     )
else
    CONTAINER_IMAGES=("${CONTAINER_IMAGE}")
fi


declare -a RESULT_MATRIX
for index in "${!MATRIX_FLAGS[@]}"; do RESULT_MATRIX[$index]=0; done


##
# Print out the command to be executed, and execute it.
##
function verbose
{
    echo "$@" >&2
    "$@"
} # verbose


##
# Run a test iteration with the environment variables
# specified.
##
function run-test
{
    local delay
    local container_id
    local pid_logs
    local reto


    function exist-container
    {
        [ "$container_id" == "" ] && return 1
        $oci container exec "${container_id}" true &>/dev/null || return 2
        return 0
    }

    function kill-container
    {
        echo "Killing container ${container_id}"
        exist-container && {
            $oci container kill "${container_id}" 1>/dev/null
            sleep 1
        }
        exist-container && {
            $oci container kill --signal KILL "${container_id}" &>/dev/null
            sleep 1
        }
    }

    reto=0
    delay="$1"
    shift 1
    deadline=$( date +%s )
    deadline=$(( deadline + delay ))
    [ "$delay" == "" ] && echo "ERROR:No timeout specified" >&2 && exit 1

    container_id="$( STDIN_STREAM="${MATRIX_INPUTS[$index]}" \
    ${oci} run  -d \
                --volume "$PWD:$PWD:z" \
                --workdir "$PWD" \
                -e FLAG_INSTALL_VSCODE="${FLAG_INSTALL_VSCODE}" \
                -e FLAG_RUN_VSCODE_AFTER_INSTALL="${FLAG_RUN_VSCODE_AFTER_INSTALL}" \
                -e FLAG_INSTALL_CRC="${FLAG_INSTALL_CRC}" \
                fedora:32 \
                /bin/bash -c "./devel/install-local-tools.sh"
    )"
    $oci container logs --follow "${container_id}" &
    pid_logs="$!"
    trap 'kill-container' SIGINT
    while [ "$( date +%s )" -lt "$deadline" ]
    do
        $oci container exec "${container_id}" true 2>/dev/null || break
        sleep 3
    done
    if [ "$( date +%s )" -ge "$deadline" ]
    then
        # Timeout
        kill ${pid_logs}
        wait ${pid_logs}
        kill-container
        return 127
    fi
    reto="$( $oci container wait "${container_id}" 2>/dev/null )"
    $oci container rm -f "${container_id}" &>/dev/null
    # shellcheck disable=SC2086
    return $reto
}


START_INDEX=${START_INDEX:-1}
END_INDEX=${END_INDEX:-40}

if command -v podman 1>/dev/null 2>/dev/null
then
    oci="podman"
elif command -v docker 1>/dev/null 2>/dev/null
then
    oci="docker"
else
    echo "ERROR:podman nor docker were found"
    exit 3
fi


###############################################################################

# Now launch the tests for all the different scenarios

###############################################################################

true > ".test-install-local-tools.report"
for image in "${CONTAINER_IMAGES[@]}"
do
    for index in "${!MATRIX_FLAGS[@]}"
    do
        # shellcheck disable=SC2086
        [ $index -lt $START_INDEX ] && continue
        # shellcheck disable=SC2086
        [ $index -gt $END_INDEX ] && continue
        # shellcheck disable=SC2206
        CURRENT_FLAGS=(${MATRIX_FLAGS[$index]})
        echo "INFO:image='${image}'; CURRENT_FLAGS[$index]=${CURRENT_FLAGS[*]}; MATRIX_EXPECTED=${MATRIX_EXPECTED[$index]}; MATRIX_INPUTS[$index]=${MATRIX_INPUTS[$index]}" >&2
        FLAG_INSTALL_VSCODE="${CURRENT_FLAGS[0]}"
        FLAG_RUN_VSCODE_AFTER_INSTALL="${CURRENT_FLAGS[1]}"
        FLAG_INSTALL_CRC="${CURRENT_FLAGS[2]}"
        run-test "${DEFAULT_TIMEOUT}"
        RESULT_MATRIX[$index]=$?
        echo "INFO:RESULT=${RESULT_MATRIX[$index]}" >&2
    done

    # Print result report
    for index in "${!MATRIX_FLAGS[@]}"
    do
        # shellcheck disable=SC2086
        [ $index -lt $START_INDEX ] && continue
        # shellcheck disable=SC2086
        [ $index -gt $END_INDEX ] && continue
        if [ ${RESULT_MATRIX[$index]} -eq ${MATRIX_EXPECTED[$index]} ]
        then
            RESULT="SUCCESS"
        else
            RESULT="FAILURE; RESULT=${RESULT_MATRIX[$index]}; EXPECTED=${MATRIX_EXPECTED[$index]}"
        fi
        echo "image='${image}'; index=$index; CURRENT_FLAGS[$index]=${CURRENT_FLAGS[*]}; MATRIX_EXPECTED=${MATRIX_EXPECTED[$index]}: ${RESULT}" >>tests/.test-install-local-tools.report
    done
done
