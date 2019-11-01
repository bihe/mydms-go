import { ReplaySubject } from 'rxjs';
import { ApplicationData } from '../models/application.data';

const ShowAmountKey = 'mydms.amount.show';

export class ApplicationState {
  private appData: ReplaySubject<ApplicationData> = new ReplaySubject();
  private searchInput: ReplaySubject<string> = new ReplaySubject();
  private progress: ReplaySubject<boolean> = new ReplaySubject();
  private showAmount: ReplaySubject<boolean> = new ReplaySubject();
  private requestReload: ReplaySubject<boolean> = new ReplaySubject();

  public getAppData(): ReplaySubject<ApplicationData> {
    return this.appData;
  }

  public setAppData(data: ApplicationData) {
    this.appData.next(data);
  }

  public getSearchInput(): ReplaySubject<string> {
    return this.searchInput;
  }

  public setSearchInput(data: string) {
    this.searchInput.next(data);
  }

  public setProgress(data: boolean) {
    this.progress.next(data);
  }

  public getProgress(): ReplaySubject<boolean> {
    return this.progress;
  }

  public getShowAmount(): ReplaySubject<boolean> {
    const str = localStorage.getItem(ShowAmountKey);
    if (!str || str === undefined) {
      this.showAmount.next(false);
      return this.showAmount;
    }
    const parsed = JSON.parse(str);
    let val = false;
    if (parsed === true || parsed === 1 || parsed === "1" || parsed === "true" || parsed === "TRUE") {
      val = true;
    }
    this.showAmount.next(val);
    return this.showAmount;
  }

  public setShowAmount(show: boolean) {
    localStorage.setItem(ShowAmountKey, JSON.stringify(show));
    this.showAmount.next(show);
  }

  public setRequestReload(data: boolean) {
    this.requestReload.next(data);
  }

  public getRequestReload(): ReplaySubject<boolean> {
    return this.requestReload;
  }
}
