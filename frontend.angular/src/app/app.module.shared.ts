import { NgModule } from '@angular/core';
// tslint:disable-next-line:max-line-length
import { MatDialogModule, MatInputModule, MatNativeDateModule, MatProgressSpinnerModule, MatSnackBarModule, MatTooltipModule } from '@angular/material';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { TagInputModule } from 'ngx-chips';
import { NgxUploaderModule } from 'ngx-uploader';
import { AppComponent } from './components/app/app.component';
import { DocumentComponent } from './components/document/document.component';
import { HomeComponent } from './components/home/home.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { ConfirmationDialogComponent } from './shared/confirmation/confirmation.component';
import { DateFormatPipe } from './shared/pipes/date.format';
import { EllipsisPipe } from './shared/pipes/ellipsis';
import { AppDataService } from './shared/services/app.data.service';
import { ApplicationState } from './shared/services/app.state';


@NgModule({
  imports: [ MatInputModule, MatSnackBarModule, MatDialogModule, MatNativeDateModule, MatProgressSpinnerModule,
    MatTooltipModule, MatChipsModule, MatCardModule, MatSlideToggleModule, MatButtonModule ],
  exports: [ MatInputModule, MatSnackBarModule, MatDialogModule, MatNativeDateModule, MatProgressSpinnerModule,
    MatTooltipModule, MatChipsModule, MatCardModule, MatSlideToggleModule, MatButtonModule ],
})
export class AppMaterialModule { }

export const sharedConfig: NgModule = {
    bootstrap: [ AppComponent ],
    declarations: [
        AppComponent,
        NavbarComponent,
        HomeComponent,
        ConfirmationDialogComponent,
        DocumentComponent,
        DateFormatPipe,
        EllipsisPipe
    ],
    imports: [
        AppMaterialModule,
        TagInputModule,
        NgxUploaderModule
    ],
    providers: [ ApplicationState, AppDataService ],
    entryComponents: [ ConfirmationDialogComponent ]
};
