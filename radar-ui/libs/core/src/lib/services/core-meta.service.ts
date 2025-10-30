import { Injectable } from '@angular/core';
import { Meta, Title } from '@angular/platform-browser';

import { PageMeta } from '../interfaces/core-meta.interface';

@Injectable({
    providedIn: 'root'
})
export class CoreMetaService {
    constructor(
        private readonly metaService: Meta,
        private readonly titleService: Title
    ) {}

    initPageMetaTags() {
        this.metaService.addTags([
            {
                name: 'description',
                content: ''
            }
        ]);
    }

    setPageMeta(meta: PageMeta) {
        this.titleService.setTitle(meta.title);
        this.metaService.updateTag({
            name: 'description',
            content: meta.description || ''
        });
    }
}
