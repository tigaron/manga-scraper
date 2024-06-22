-- DropTable
DROP TABLE IF EXISTS `SeriesListData`;
DROP TABLE IF EXISTS `SeriesDetailData`;
DROP TABLE IF EXISTS `ChapterListData`;
DROP TABLE IF EXISTS `ChapterDetailData`;

-- AlterTable
ALTER TABLE `ScrapeRequest` TTL = `createdAt` + INTERVAL 1 WEEK;
