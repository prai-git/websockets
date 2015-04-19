<?php

if (count($argv) < 2) {
    echo "usage: php ./zmq.php <message>\n";
}

$message = $argv[1];

$context = new \ZMQContext();
$socket = $context->getSocket(\ZMQ::SOCKET_PUB, 'message');
$socket->bind("tcp://127.0.0.1:5563");

echo "Sending: ".$message."\n";

sleep(1);

$socket->sendmulti(array("message", $message));