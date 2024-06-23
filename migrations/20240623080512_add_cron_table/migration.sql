-- CreateTable
CREATE TABLE `CronJob` (
    `id` VARCHAR(191) NOT NULL,
    `name` TEXT NOT NULL,
    `crontab` TEXT NOT NULL,
    `tags` TEXT NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `CronJobStatus` (
    `id` VARCHAR(191) NOT NULL,
    `jobId` VARCHAR(191) NOT NULL,
    `status` TEXT NOT NULL,
    `message` TEXT NOT NULL,
    `duration` DOUBLE NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `cronJobIndex`(`jobId`),
    PRIMARY KEY (`id`)
) TTL = `createdAt` + INTERVAL 1 MONTH DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `CronJobStatus`
ADD CONSTRAINT `CronJobStatus_jobId_fkey`
FOREIGN KEY (`jobId`)
REFERENCES `CronJob` (`id`)
ON DELETE CASCADE ON UPDATE CASCADE;
