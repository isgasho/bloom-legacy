import 'dart:async';

import 'package:bloom/kernel/blocs/bloc_provider.dart';
import 'package:contacts_service/contacts_service.dart';
import 'package:flutter/material.dart';

class ContactBloc extends BlocBase {
  ContactBloc({@required Contact contact}) {
    _contact = contact;
  }

  Contact _contact;

  final StreamController<Contact> _contactController =
      StreamController<Contact>.broadcast();
  // StreamSink<Contact> get _inContact => _contactController.sink;
  Stream<Contact> get contactOut => _contactController.stream;

  final StreamController<Contact> _contactDeletedController =
      StreamController<Contact>.broadcast();
  StreamSink<Contact> get _inDeleted => _contactDeletedController.sink;
  Stream<Contact> get deleted => _contactDeletedController.stream;

  Contact get contact => _contact;

  @override
  void dispose() {
    _contactController.close();
  }

  Future<void> delete() async {
    await ContactsService.deleteContact(_contact);
    _inDeleted.add(_contact);
  }
}
