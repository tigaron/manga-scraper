-- AddChaptersCount
ALTER TABLE `Series` ADD COLUMN `chaptersCount` INT NOT NULL DEFAULT 0;

-- AddLatestChapter
ALTER TABLE `Series` ADD COLUMN `latestChapter` VARCHAR(191) NOT NULL DEFAULT '';
