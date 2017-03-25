import { Injectable } from '@angular/core';

@Injectable()
export class WsService {
  private conn: WebSocket;
  private queue: string[] = [];

  constructor() {
    this.conn = new WebSocket('ws://localhost:8080/ws');
    this.conn.addEventListener('open', () => {
      this.queue.forEach(msg => this.conn.send(msg));
      this.queue = [];
    });
  }

  send(msg: string) {
    if (this.conn.readyState === WebSocket.OPEN) {
      this.conn.send(msg);
    } else {
      this.queue.push(msg);
    }
  }

  onmessage(cb: (string) => void) {
    this.conn.addEventListener('message', e => cb(JSON.parse(e.data)));
  }
}
