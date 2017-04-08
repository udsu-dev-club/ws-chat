import { Injectable, EventEmitter } from '@angular/core';
import { Observable } from 'rxjs/Rx';

@Injectable()
export class WsService {
    private _conn: WebSocket;
    private _queue: string[] = [];
    private _logined = false;
    private _id = 1;
    
    public response: EventEmitter<any>;
    
    get logined(): boolean {
        return this._logined;
    }
    
    constructor() {
        this._conn = new WebSocket('ws://localhost:8080/ws');
        this._conn.addEventListener('open', () => {
            this._queue.forEach(msg => {
                this._conn.send(msg)
            });
            this._queue = [];
        });
        this.response = new EventEmitter<any>();
        this._conn.addEventListener('message', e => {
            const data = JSON.parse(e.data);
            this.response.emit(data);
        });
        this.response
            .subscribe(msg => {
                if (msg.cmd === 'LOGIN' && msg.id > 0) {
                    this._logined = true;
                }
            })
        ;
        this.response
            .subscribe(msg => {
                if (msg.cmd === 'LOGOUT' && msg.id > 0 && !msg.error) {
                    this._logined = false;
                }
            })
        ;
    }
    
    send(cmd: string, data: any): Promise<any> {
        const req = {
            id: this._id++,
            cmd: cmd,
            data: data
        }
        const res = new Promise<any>(resolve => {
            const sub = this.response
                .subscribe(msg => {
                    if (msg.id === req.id) {
                        sub.unsubscribe();
                        resolve(msg);
                    }
                });
        });

        const sreq = JSON.stringify(req);

        if (this._conn.readyState === WebSocket.OPEN) {
            this._conn.send(sreq);
        } else {
            this._queue.push(sreq);
        }

        return res;
    }
}
