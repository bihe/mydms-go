import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material';
import { Router } from '@angular/router';
import * as moment from 'moment';
import { debounceTime } from 'rxjs/operators';
import { Document } from '../../shared/models/document.model';
import { BackendService } from '../../shared/services/backend.service';
import { ApplicationState } from '../../shared/services/state.service';
import { MessageUtils } from '../../shared/utils/message.utils';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  documents: Array<Document> = new Array<Document>();
  totalEntries = 0;
  shownResults = 0;
  showAmount = false;

  readonly InitialPageSize: number = 20;
  private pagedDocuments: Array<Document> = null;
  private searchString: string = null;

  constructor(
    private state: ApplicationState,
    private service: BackendService,
    private router: Router,
    private snackBar: MatSnackBar) {

    this.state.getSearchInput().pipe(
      debounceTime(300))
      .subscribe(x => {
        console.log('Search for: ' + x);
        this.searchString = x;
        this.documents = [];
        this.searchDocuments(x, 0);
    });

    this.state.getShowAmount().subscribe(
      x => {
        this.showAmount = x;
      }
    );

    this.state.getRequestReload().subscribe(
      reload => {
        if (reload) {
          this.documents = new Array<Document>();
          this.searchDocuments(null, 0);
        }
      }
    );
  }

  ngOnInit() {
    this.searchDocuments(null, 0);
  }

  showMoreResults() {
    this.searchDocuments(this.searchString, this.shownResults);
  }

  clearSearch() {
    this.state.setSearchInput('');
  }

  addDocument() {
    this.router.navigate(['/document/-1']);
  }

  editDocument(doc: Document) {
    this.router.navigate(['/document/' + doc.id]);
  }

  searchDocuments(title: string, skipEntries: number) {
    this.state.setProgress(true);
    this.service.searchDocuments(title, this.InitialPageSize, skipEntries)
      .subscribe(
        result => {
          const returnedResults = (result.documents) ? result.documents.length : 0;
          this.totalEntries = result.totalEntries;
          this.shownResults = skipEntries + returnedResults;

          const doucmentResult = result.documents;
          console.log('Result from search: ' + returnedResults);
          if (doucmentResult) {
            this.pagedDocuments = new Array<Document>();
            doucmentResult.forEach(a => {
              const doc = new Document();
              doc.title = a.title;
              doc.created = a.created;
              doc.modified = a.modified;
              if (this.showAmount) {
                doc.amount = a.amount;
              }
              doc.fileName = a.fileName;
              doc.encodedFilename = btoa(encodeURI(a.fileName));
              doc.previewLink = a.previewLink;
              doc.id = a.id;
              doc.tags = a.tags;
              doc.senders = a.senders;
              doc.dateHuman = moment(doc.lastDate).fromNow();

              this.pagedDocuments.push(doc);
            });
            this.documents = this.documents.concat(this.pagedDocuments);
            this.documents = arrayUnique(this.documents);
          }
          this.state.setProgress(false);
        },
        error => {
          this.state.setProgress(false);
          new MessageUtils().showError(this.snackBar, error);
        }
      );
  }
}

function arrayUnique(array: any) {
    const a = array.concat();
    for (let i = 0; i < a.length; ++i) {
        for (let j = i + 1; j < a.length; ++j) {
            if (a[i].id === a[j].id) {
              a.splice(j--, 1);
            }
        }
    }
    return a;
}
