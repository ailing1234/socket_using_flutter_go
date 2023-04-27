import 'dart:convert';

import 'package:web_socket_channel/io.dart';
import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:web_socket_channel/io.dart'; // for non-web platforms
import 'package:web_socket_channel/html.dart'; // for web platform

void main() {
  final channel = kIsWeb
      ? HtmlWebSocketChannel.connect("ws://localhost:8080/ws")
      : IOWebSocketChannel.connect("ws://localhost:8080/ws");

  runApp(MyApp(channel: channel));
}

class MyApp extends StatefulWidget {
  final channel;

  MyApp({this.channel});

  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  String _text = 'Connecting...';

  @override
  void initState() {
    super.initState();

    widget.channel.stream.listen((message) {
      print("Message: listen" + message.toString());
      var data = jsonDecode(message);
      print("message: " + data['message']);

      var text;
      if (data['channel'] == 'channel1') {
        // handle chat message
        text = data['channel'];
        var sender = data['message'];
        print('BBL : $sender: $text');
      } else if (data['channel'] == 'notification') {
        // handle notification message
        text = data['message'];
        print('Notification: $text');
      } else {
        text = data['message'];
        print('A: $text');
      }

      setState(() {
        _text = text;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: const Text('WebSocket Example'),
        ),
        body: Center(
          child: Text(_text),
        ),
      ),
    );
  }

  @override
  void dispose() {
    widget.channel.sink.close();
    super.dispose();
  }
}
