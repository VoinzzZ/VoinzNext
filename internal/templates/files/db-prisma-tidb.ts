import { PrismaTiDBCloud } from "@tidbcloud/prisma-adapter";
import { connect } from "@tidbcloud/serverless";
import { PrismaClient } from "@prisma/client";

const connectionString = process.env.DATABASE_URL!;
const adapter = new PrismaTiDBCloud({ url: connectionString });
const globalForPrisma = globalThis as unknown as {
  prisma: PrismaClient | undefined;
};

export const db = globalForPrisma.prisma ?? new PrismaClient({ adapter });

if (process.env.NODE_ENV !== "production") globalForPrisma.prisma = db;
