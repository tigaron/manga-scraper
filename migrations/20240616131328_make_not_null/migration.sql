-- UpdateTable
UPDATE `Series`
SET `thumbnailUrl` = COALESCE(`thumbnailUrl`, '');

ALTER TABLE `Series`
MODIFY COLUMN `thumbnailUrl` TEXT NOT NULL;

-- UpdateTable
UPDATE `Series`
SET `synopsis` = COALESCE(`synopsis`, '');

ALTER TABLE `Series`
MODIFY COLUMN `synopsis` TEXT NOT NULL;

-- UpdateTable
UPDATE `Series`
SET `genres` = COALESCE(`genres`, '[]');

ALTER TABLE `Series`
MODIFY COLUMN `genres` JSON NOT NULL;

-- UpdateTable
UPDATE `Chapter`
SET `fullTitle` = COALESCE(`fullTitle`, '');

ALTER TABLE `Chapter`
MODIFY COLUMN `fullTitle` TEXT NOT NULL;

-- UpdateTable
UPDATE `Chapter`
SET `sourcePath` = COALESCE(`sourcePath`, '');

ALTER TABLE `Chapter`
MODIFY COLUMN `sourcePath` TEXT NOT NULL;

-- UpdateTable
UPDATE `Chapter`
SET `nextSlug` = COALESCE(`nextSlug`, '');

ALTER TABLE `Chapter`
MODIFY COLUMN `nextSlug` TEXT NOT NULL;

-- UpdateTable
UPDATE `Chapter`
SET `nextPath` = COALESCE(`nextPath`, '');

ALTER TABLE `Chapter`
MODIFY COLUMN `nextPath` TEXT NOT NULL;

-- UpdateTable
UPDATE `Chapter`
SET `prevSlug` = COALESCE(`prevSlug`, '');

ALTER TABLE `Chapter`
MODIFY COLUMN `prevSlug` TEXT NOT NULL;

-- UpdateTable
UPDATE `Chapter`
SET `prevPath` = COALESCE(`prevPath`, '');

ALTER TABLE `Chapter`
MODIFY COLUMN `prevPath` TEXT NOT NULL;

-- UpdateTable
UPDATE `Chapter`
SET `contentPaths` = COALESCE(`contentPaths`, '[]');

ALTER TABLE `Chapter`
MODIFY COLUMN `contentPaths` JSON NOT NULL;

-- UpdateTable
UPDATE `ScrapeRequest`
SET `series` = COALESCE(`series`, '');

ALTER TABLE `ScrapeRequest`
MODIFY COLUMN `series` TEXT NOT NULL;

-- UpdateTable
UPDATE `ScrapeRequest`
SET `chapter` = COALESCE(`chapter`, '');

ALTER TABLE `ScrapeRequest`
MODIFY COLUMN `chapter` TEXT NOT NULL;

-- UpdateTable
UPDATE `ScrapeRequest`
SET `status` = COALESCE(`status`, '');

ALTER TABLE `ScrapeRequest`
MODIFY COLUMN `status` VARCHAR(191) NOT NULL;

-- UpdateTable
UPDATE `ScrapeRequest`
SET `retries` = COALESCE(`retries`, 0);

ALTER TABLE `ScrapeRequest`
MODIFY COLUMN `retries` INTEGER NOT NULL;

-- UpdateTable
UPDATE `ScrapeRequest`
SET `totalTime` = COALESCE(`totalTime`, 0);

ALTER TABLE `ScrapeRequest`
MODIFY COLUMN `totalTime` DOUBLE NOT NULL;

-- UpdateTable
UPDATE `ScrapeRequest`
SET `error` = COALESCE(`error`, false);

ALTER TABLE `ScrapeRequest`
MODIFY COLUMN `error` BOOLEAN NOT NULL;

-- UpdateTable
UPDATE `ScrapeRequest`
SET `message` = COALESCE(`message`, '');

ALTER TABLE `ScrapeRequest`
MODIFY COLUMN `message` TEXT NOT NULL;

-- UpdateTable
UPDATE `ChapterDetailData`
SET `nextSlug` = COALESCE(`nextSlug`, '');

ALTER TABLE `ChapterDetailData`
MODIFY COLUMN `nextSlug` TEXT NOT NULL;

-- UpdateTable
UPDATE `ChapterDetailData`
SET `nextPath` = COALESCE(`nextPath`, '');

ALTER TABLE `ChapterDetailData`
MODIFY COLUMN `nextPath` TEXT NOT NULL;

-- UpdateTable
UPDATE `ChapterDetailData`
SET `prevSlug` = COALESCE(`prevSlug`, '');

ALTER TABLE `ChapterDetailData`
MODIFY COLUMN `prevSlug` TEXT NOT NULL;

-- UpdateTable
UPDATE `ChapterDetailData`
SET `prevPath` = COALESCE(`prevPath`, '');

ALTER TABLE `ChapterDetailData`
MODIFY COLUMN `prevPath` TEXT NOT NULL;

-- UpdateTable
UPDATE `ChapterDetailData`
SET `contentPaths` = COALESCE(`contentPaths`, '[]');

ALTER TABLE `ChapterDetailData`
MODIFY COLUMN `contentPaths` JSON NOT NULL;
