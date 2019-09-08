import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogConfig } from '@angular/material';

@Component({
  selector: 'app-confirmation',
  templateUrl: './confirmation.component.html',
  styleUrls: ['./confirmation.component.css']
})
export class ConfirmationDialogComponent implements OnInit {

    name = '';

    constructor(
        private dialog: MatDialogRef<ConfirmationDialogComponent>,
        @Inject(MAT_DIALOG_DATA) private dialogData: any) {
        this.name = dialogData.name;
    }

    public static getDialogConfig(name: string): MatDialogConfig {
        const dialogConfig: MatDialogConfig = {
            disableClose: false,
            width: '440px',
            height: '240px',
            data: {
                name: name
            }
        };
        return dialogConfig;
    }

    public onConfirmDelete() {
        this.dialog.close(true);
    }

    ngOnInit() {
    }
}
