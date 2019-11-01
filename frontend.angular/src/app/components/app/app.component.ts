import { Component } from '@angular/core';
import { ApplicationState } from '../../shared/services/state.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent {

  public version: string;
  public displayName: string;
  public year: number;
  public build: string;
  public runtimeVersion: string;

  constructor(private state: ApplicationState) {
    this.state.getAppData().subscribe(x => {
      this.displayName = x.appInfo.userInfo.displayName;
      this.version = x.appInfo.versionInfo.version;
      this.build = x.appInfo.versionInfo.buildNumber;
      this.runtimeVersion = x.appInfo.versionInfo.runtimeVersion;
    });
    this.year = new Date().getFullYear();
  }
}
