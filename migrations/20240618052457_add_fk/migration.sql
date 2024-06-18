-- AddForeignKey
ALTER TABLE `Series` ADD CONSTRAINT `Series_providerSlug_fkey` FOREIGN KEY (`providerSlug`) REFERENCES `Provider`(`slug`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Chapter` ADD CONSTRAINT `Chapter_providerSlug_fkey` FOREIGN KEY (`providerSlug`) REFERENCES `Provider`(`slug`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Chapter` ADD CONSTRAINT `Chapter_providerSlug_seriesSlug_fkey` FOREIGN KEY (`providerSlug`, `seriesSlug`) REFERENCES `Series`(`providerSlug`, `slug`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `SeriesListData` ADD CONSTRAINT `SeriesListData_requestId_fkey` FOREIGN KEY (`requestId`) REFERENCES `ScrapeRequest`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `SeriesDetailData` ADD CONSTRAINT `SeriesDetailData_requestId_fkey` FOREIGN KEY (`requestId`) REFERENCES `ScrapeRequest`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ChapterListData` ADD CONSTRAINT `ChapterListData_requestId_fkey` FOREIGN KEY (`requestId`) REFERENCES `ScrapeRequest`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ChapterDetailData` ADD CONSTRAINT `ChapterDetailData_requestId_fkey` FOREIGN KEY (`requestId`) REFERENCES `ScrapeRequest`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;
