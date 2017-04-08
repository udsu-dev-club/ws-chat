import { WsService } from './../ws.service';
import { IMessage } from './chat.component.stub';
import { Component, OnInit, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { Subscriber } from 'rxjs/Rx';

@Component({
    selector: 'app-chat',
    templateUrl: './chat.component.html',
    styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit, OnDestroy {
    private ws: WsService;
    private body: string;
    private messages: IMessage[] = [];
    private sub: Subscriber<any>;

    @ViewChild('username') username: ElementRef;
    
    constructor(ws: WsService) {
        this.ws = ws;
    }
    
    ngOnInit() {
        this.sub = this.ws.response.subscribe(msg => {
            // console.log(msg);
        });
    }

    ngOnDestroy() {
        this.sub.unsubscribe();
    }

    login() {
        this.ws.send('LOGIN', {username: this.username.nativeElement.value});
    }
    
    // send(event: KeyboardEvent) {
    //     if (event.keyCode === 13) {
    //         this.ws.send(JSON.stringify(this.body));
    //         this.body = '';
    //     }
    // }
}
