import { Component, OnInit } from '@angular/core';
import { MatSlideToggleChange, MatSnackBar } from '@angular/material';
import { DomSanitizer } from '@angular/platform-browser';
import { Router } from '@angular/router';
import { ApplicationData } from '../../shared/models/application.data';
import { BackendService } from '../../shared/services/backend.service';
import { ApplicationState } from '../../shared/services/state.service';
import { MessageUtils } from '../../shared/utils/message.utils';

@Component({
  selector: 'app-nav-bar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {

  menuVisible = false;
  showProgress = false;
  searchText = '';
  showAmount = false;

  public A: ApplicationData;

  constructor(
    private service: BackendService,
    private state: ApplicationState,
    private snackBar: MatSnackBar,
    private sanitizer: DomSanitizer,
    private router: Router) {
  }

  ngOnInit() {
    this.service.getApplicationInfo()
      .subscribe(
        data => {
          this.A = new ApplicationData();
          this.A.appInfo = data;
          this.state.setAppData(this.A);
        },
        error => {
          new MessageUtils().showError(this.snackBar, error);
        }
      );

    this.state.getSearchInput()
      .subscribe(
        data => {
          this.searchText = data;
        },
        error => {
          new MessageUtils().showError(this.snackBar, error);
        }
      );

    this.state.getProgress()
      .subscribe(
        data => {
          this.showProgress = data;
        },
        error => {
          new MessageUtils().showError(this.snackBar, error);
        }
      );

    this.state.getShowAmount()
    .subscribe(
      x => {
        this.showAmount = x;
      }
    );

  }

  onSearch(searchText: string) {
    this.state.setSearchInput(searchText);
  }

  toggleMenu(visible: boolean) {
    this.menuVisible = visible;
  }

  menuTransform() {
    if (this.menuVisible) {
      return this.sanitizer.bypassSecurityTrustStyle('translateX(0)');
    } else {
      return this.sanitizer.bypassSecurityTrustStyle('translateX(-110%)');
    }
  }

  navigateTo(destination: string) {
    this.toggleMenu(false);
    this.router.navigate([destination]);
  }

  showAmountToggle(event: MatSlideToggleChange) {
    console.log('Change showAmount to ' + event.checked);
    this.state.setShowAmount(event.checked);
    this.state.setRequestReload(true);
  }
}
