-- CreateTable
CREATE TABLE `Provider` (
    `id` VARCHAR(191) NOT NULL,
    `slug` VARCHAR(191) NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `scheme` VARCHAR(191) NOT NULL,
    `host` TEXT NOT NULL,
    `listPath` TEXT NOT NULL,
    `isActive` BOOLEAN NOT NULL DEFAULT false,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `Provider_slug_key`(`slug`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Series` (
    `id` VARCHAR(191) NOT NULL,
    `slug` VARCHAR(191) NOT NULL,
    `title` TEXT NOT NULL,
    `sourcePath` TEXT NOT NULL,
    `thumbnailUrl` TEXT NULL,
    `synopsis` TEXT NULL,
    `genres` JSON NULL,
    `providerSlug` VARCHAR(191) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `providerIndex`(`providerSlug`),
    UNIQUE INDEX `Series_providerSlug_slug_key`(`providerSlug`, `slug`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Chapter` (
    `id` VARCHAR(191) NOT NULL,
    `slug` VARCHAR(191) NOT NULL,
    `number` DOUBLE NOT NULL,
    `shortTitle` TEXT NOT NULL,
    `sourceHref` TEXT NOT NULL,
    `fullTitle` TEXT NULL,
    `sourcePath` TEXT NULL,
    `nextSlug` TEXT NULL,
    `nextPath` TEXT NULL,
    `prevSlug` TEXT NULL,
    `prevPath` TEXT NULL,
    `contentPaths` JSON NULL,
    `providerSlug` VARCHAR(191) NOT NULL,
    `seriesSlug` VARCHAR(191) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `providerIndex`(`providerSlug`),
    INDEX `seriesIndex`(`providerSlug`, `seriesSlug`),
    UNIQUE INDEX `Chapter_providerSlug_seriesSlug_slug_key`(`providerSlug`, `seriesSlug`, `slug`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ScrapeRequest` (
    `id` VARCHAR(191) NOT NULL,
    `type` ENUM('SERIES_LIST', 'SERIES_DETAIL', 'CHAPTER_LIST', 'CHAPTER_DETAIL') NOT NULL,
    `baseUrl` TEXT NOT NULL,
    `requestPath` TEXT NOT NULL,
    `provider` TEXT NOT NULL,
    `series` TEXT NULL,
    `chapter` TEXT NULL,
    `status` VARCHAR(191) NULL,
    `retries` INTEGER NULL DEFAULT 0,
    `totalTime` DOUBLE NULL,
    `error` BOOLEAN NULL DEFAULT false,
    `message` TEXT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `typeIndex`(`type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `SeriesListData` (
    `id` VARCHAR(191) NOT NULL,
    `title` TEXT NOT NULL,
    `slug` TEXT NOT NULL,
    `sourcePath` TEXT NOT NULL,
    `requestId` VARCHAR(191) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `requestIndex`(`requestId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `SeriesDetailData` (
    `id` VARCHAR(191) NOT NULL,
    `thumbnailUrl` TEXT NOT NULL,
    `synopsis` TEXT NOT NULL,
    `genres` JSON NOT NULL,
    `requestId` VARCHAR(191) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `requestIndex`(`requestId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ChapterListData` (
    `id` VARCHAR(191) NOT NULL,
    `shortTitle` TEXT NOT NULL,
    `slug` TEXT NOT NULL,
    `number` DOUBLE NOT NULL,
    `href` TEXT NOT NULL,
    `requestId` VARCHAR(191) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `requestIndex`(`requestId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ChapterDetailData` (
    `id` VARCHAR(191) NOT NULL,
    `fullTitle` TEXT NOT NULL,
    `sourcePath` TEXT NOT NULL,
    `nextSlug` TEXT NULL,
    `nextPath` TEXT NULL,
    `prevSlug` TEXT NULL,
    `prevPath` TEXT NULL,
    `contentPaths` JSON NOT NULL,
    `requestId` VARCHAR(191) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `requestIndex`(`requestId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
