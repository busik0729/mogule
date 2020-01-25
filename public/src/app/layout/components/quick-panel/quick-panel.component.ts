import { Component, ViewEncapsulation } from '@angular/core';
import { IMessage } from "../../../app.component";
import { WS } from "../../../websocket/websocket.events";
import { WebsocketService } from "../../../websocket";

@Component({
    selector     : 'quick-panel',
    templateUrl  : './quick-panel.component.html',
    styleUrls    : ['./quick-panel.component.scss'],
    encapsulation: ViewEncapsulation.None
})
export class QuickPanelComponent
{
    date: Date;
    events: any[];
    notes: any[];
    settings: any;

    /**
     * Constructor
     */
    constructor(private _socket: WebsocketService)
    {
        // Set the defaults
        this.date = new Date();
        this.settings = {
            notify: true,
            cloud : false,
            retro : true
        };
        this.events = [];


        this.registerWsListeners()
    }

    registerWsListeners() {
        this._socket.on<IMessage[]>(WS.ON.BOARD.DEADLINE).subscribe((res: any) => {
            console.log(res);
            const d = Date.parse(res.data.due);
            const dd = new Date(d).toLocaleDateString();
            this.events = this.events.filter((event) => event.id != res.data.id);
            this.events.unshift({
                id: res.id + "",
                title: "Дэдлайн проекта",
                detail: res.data.name + " " + dd,
            })
        });

        this._socket.on<IMessage[]>(WS.ON.CARD.DEADLINE).subscribe((res: any) => {
            console.log(res);
            const d = Date.parse(res.data.due);
            const dd = new Date(d).toLocaleDateString();
            this.events = this.events.filter((event) => event.id != res.data.id);
            this.events.unshift({
                id: res.id + "",
                title: "Дэдлайн задания",
                detail: res.data.name + " " + dd,
            })
        });
    }

    closeEvent(event: any) {
        this._socket.send(WS.SEND.READ_EVENT, {id: event.id});

        const index = this.events.indexOf(event);

        this.events.splice(index, 1);
    }
}
