import { WsService } from './../ws.service';
import { IMessage } from './chat.component.stub';
import { Component, OnInit } from '@angular/core';

@Component({
    selector: 'app-chat',
    templateUrl: './chat.component.html',
    styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {
    private ws: WsService;
    private body: string;
    private messages: IMessage[] = [];
    
    constructor(ws: WsService) {
        this.ws = ws;
    }
    
    ngOnInit() {
        this.ws.onmessage(msg => {
            this.messages.push({
                timestamp: (new Date()).toDateString(),
                author: '@' + msg.slice(0, 16),
                body: msg.slice(16),
            });
        });
    }
    
    send(event: KeyboardEvent) {
        if (event.keyCode === 13) {
            this.ws.send(JSON.stringify(this.body));
            this.body = '';
        }
    }
}
