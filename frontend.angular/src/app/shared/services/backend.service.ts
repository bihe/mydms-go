import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError, distinctUntilChanged, timeout } from 'rxjs/operators';
import { AppInfo } from '../models/app.info';
import { Document } from '../models/document.model';
import { DocumentResult } from '../models/document.result.model';
import { Result } from '../models/result.model';
import { SearchResult } from '../models/result.search';
import { BaseService } from './_base.service';


@Injectable()
export class BackendService extends BaseService {
  private readonly SEARCH_DOCUMENTS: string = '/api/v1/documents/search';
  private readonly APP_INFO_URL: string = '/api/v1/appinfo';
  private readonly SAVE_DOCUMENTS_URL: string = '/api/v1/documents/';
  private readonly LOAD_DOCUMENT_URL: string = '/api/v1/documents/%ID%';

  private readonly SEARCH_SENDERS_URL: string = '/api/v1/documents/senders/search';
  private readonly SEARCH_TAGS_URL: string = '/api/v1/documents/tags/search';

  constructor (private http: HttpClient) {
    super();
  }

  getApplicationInfo(): Observable<AppInfo> {
    return this.http.get<AppInfo>(this.APP_INFO_URL, this.RequestOptions)
      .pipe(
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );
  }

  searchDocuments(title: string, pageSize: number, skipEntries: number): Observable<DocumentResult> {
    const searchUrl = 'title=%TITLE%&limit=%LIMIT%&skip=%SKIP%';
    if (title && title !== '') {
      title = encodeURIComponent(title);
    }
    let url = this.SEARCH_DOCUMENTS + '?' + searchUrl.replace('%TITLE%', title || '');
    if (!pageSize) {
      url = url.replace('%LIMIT%', '');
    } else {
      url = url.replace('%LIMIT%', pageSize.toString());
    }
    if (!skipEntries) {
      url = url.replace('%SKIP%', '');
    } else {
      url = url.replace('%SKIP%', skipEntries.toString());
    }

    return this.http.get<DocumentResult>(url, this.RequestOptions)
      .pipe(
        distinctUntilChanged(),
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );
  }

  searchSenders(name: string): Observable<SearchResult> {
    const searchUrl = 'name=%NAME%';
    const url = this.SEARCH_SENDERS_URL + '?' + searchUrl.replace('%NAME%', name || '');

    return this.http.get<SearchResult>(url, this.RequestOptions)
      .pipe(
        distinctUntilChanged(),
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );
  }

  searchTags(name: string): Observable<SearchResult> {
    const searchUrl = 'name=%NAME%';
    const url = this.SEARCH_TAGS_URL + '?' + searchUrl.replace('%NAME%', name || '');

    return this.http.get<SearchResult>(url, this.RequestOptions)
      .pipe(
        distinctUntilChanged(),
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );
  }

  saveDocument(document: Document): Observable<Result> {
    return this.http.post<Result>(this.SAVE_DOCUMENTS_URL, JSON.stringify(document), this.RequestOptions)
      .pipe(
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );

  }

  getDocument(id: string): Observable<Document> {
    const url = this.LOAD_DOCUMENT_URL.replace('%ID%', id || '-1');

    return this.http.get<Document>(url, this.RequestOptions)
      .pipe(
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );
  }

  deleteDocument(id: string): Observable<Result> {
    const url = this.LOAD_DOCUMENT_URL.replace('%ID%', id || '-1');

    return this.http.delete<Result>(url, this.RequestOptions)
      .pipe(
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );
  }
}
