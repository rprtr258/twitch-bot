import logging
import subprocess

import commands.utils
import config
import utils


@commands.utils.if_empty_message("Provide pyth program: pyth.readthedocs.io")
def pyth(conf: config.Config, message_record: utils.MessageRecord):
    if message_record.username == "passionsausage":
        return "PassionSausage наказан за абьюз системы AAUGH"
    with open("c", "w", encoding="utf-8") as fd:
        fd.write(message_record.message)
    ps = subprocess.run(['echo', '6'], check=True, capture_output=True)
    # TODO: stop on timeout (5-10 seconds)
    processNames = subprocess.run(["python3", "pyth/pyth.py", "c"], input=ps.stdout, capture_output=True)
    out = processNames.stdout.decode('utf-8').replace('\r', '').replace('\n', ' ')
    logging.info("PYTH PROGRAM: %s OUTPUT: %s", repr(message_record.message), repr(out))
    if any(cmd in out for cmd in conf.pyth_config.ban_words):
        return f"@{message_record.username}, как особо опасный хакер ты приговариваешься к расстрелу! PauseFish"
    return out
