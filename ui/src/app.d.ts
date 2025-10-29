declare global {
	namespace App {
		interface Instance {
			id: string;
			name: string;
			port: number;
			status: string;
			primaryHostname: string;
			createdAt: string;
			deletedAt: string;
		}
	}
}

export {};
