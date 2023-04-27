import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:web_socket_channel/io.dart';

void main() {
  return runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter WebSocket Demo',
      theme: ThemeData(primarySwatch: Colors.blue),
      home: MyHomePage(),
    );
  }
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final channel = IOWebSocketChannel.connect('ws://localhost:8080/ws');

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Flutter WebSocket Demo'),
      ),
      body: StreamBuilder(
        stream: channel.stream,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            final message = json.decode(snapshot.data.toString());
            return Text('Received message: ${message['message']}');
          } else {
            return const Text('Waiting for messages...');
          }
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          final message = json.encode(
              {'channel': 'channel1', 'message': 'Hello server on channel1!'});
          channel.sink.add(message);
        },
        tooltip: 'Send Message',
        child: const Icon(Icons.send),
      ),
    );
  }

  @override
  void dispose() {
    channel.sink.close();
    super.dispose();
  }
}
