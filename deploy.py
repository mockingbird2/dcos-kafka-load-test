"""deploy.py

Usage:
    deploy.py [options]

Arguments:

Options:
    --dcos-username <user>        DC/OS username
    --dcos-password <pass>        DC/OS password
    --script-cpus <n>             number of CPUs to use to run the script [default: 4]
    --script-mem <n>              amount of memory (mb) to use to run the script [default: 4096]
    --brokers <s>                 Comma delimited list of Kafka brokers [default: localhost:9092]
    --compression <s>             Message compression: none, gzip, snappy [default: none]
    --creators <n>                Number of message creators [default: 1]
    --duration <n>                Duration of test in secs [default: 10]
    --event-buffer-size <n>       Overall buffered events in produceer [default: 256]
    --message-batch-size <n>      Messages per batch [default: 500]
    --message-size <n>            Message size (bytes) [default: 300]
    --produce-rate <n>            Global write rate limit (messages/sec) [default: 1000]
    --required-acks <s>           RequiredAcks config: none, local, all [default: local]
    --topic <s>                   Kafka topic which messages are send to [default: topic_test]
    --workers <n>                 Number of workers [default: 1]

"""

import os
import tempfile
import json
import subprocess
from docopt import docopt


def main(dcos_username, dcos_password, script_cpus, script_mem, topic,
         msg_size, batch_size, compression, acks, msg_rate, buffer_size,
         duration, workers, creators, brokers):
    app_defn = {
        "id": 'dcos-kafka-load-test',
        "container": {
            "type": "DOCKER",
            "docker": {
                "image": "fpaul/kafka-load-test",
            }
        },
        "cpus": script_cpus,
        "mem": script_mem,
        "disk": 1024,
        "env": {
            "DCOS_UID": dcos_username,
            "DCOS_PASSWORD": dcos_password,
            "TOPIC": topic,
            "MESSAGE_SIZE": msg_size,
            "BATCH_SIZE": batch_size,
            "COMPRESSION": compression,
            "ACKS": acks,
            "MESSAGE_RATE": msg_rate,
            "BUFFER_SIZE": buffer_size,
            "DURATION": duration,
            "WORKERS": workers,
            "CREATORS": creators,
            "BROKERS": brokers
        },
    }
    print(app_defn)
    install_app(app_defn)


def install_app(app_definition):
    def decode(output):
        return output.decode('utf-8').strip()
    app_name = app_definition["id"]

    with tempfile.TemporaryDirectory() as d:
        app_def_file = "{}.json".format(app_name.replace('/', '__'))

        app_def_path = os.path.join(d, app_def_file)

        with open(app_def_path, "w") as f:
            json.dump(app_definition, f)

        cmd = "dcos marathon app add {}".format(app_def_path)
        result = subprocess.run([cmd], shell=True, stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE)
        if result.stdout:
            print('Stdout: ', decode(result.stdout))
        if result.stderr:
            print('Stderr: ', decode(result.stderr))


if __name__ == '__main__':
    args = docopt(__doc__)

    main(dcos_username=args['--dcos-username'],
         dcos_password=args['--dcos-password'],
         script_cpus=int(args['--script-cpus']),
         script_mem=int(args['--script-mem']),
         brokers=args['--brokers'],
         compression=args['--compression'],
         creators=args['--creators'],
         duration=args['--duration'],
         buffer_size=args['--event-buffer-size'],
         batch_size=args['--message-batch-size'],
         msg_size=args['--message-size'],
         msg_rate=args['--produce-rate'],
         acks=args['--required-acks'],
         topic=args['--topic'],
         workers=args['--workers'])
