

export default async function log(info: unknown, { err } = { err: false }) {
	try {
		const { LogPrintln } = await import("$lib/bindings/App");
		await LogPrintln(`Browser - ${JSON.stringify(info)}\n`);
  } catch {/* */} finally {
		if (err) console.error(err);
		console.trace(info);
	}
}
