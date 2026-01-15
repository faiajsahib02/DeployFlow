import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { projectApi } from '@/services/api';
import type { Project } from '@/services/api';
import { Rocket, Plus, Loader2 } from 'lucide-react';

interface DashboardProps {
  onProjectCreated: (project: Project) => void;
}

export function Dashboard({ onProjectCreated }: DashboardProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [projectName, setProjectName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreateProject = async () => {
    if (!projectName.trim()) {
      setError('Project name is required');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const project = await projectApi.create(projectName.trim());
      setIsOpen(false);
      setProjectName('');
      onProjectCreated(project);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to create project. Is the backend running?');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-8">
      <div className="text-center space-y-8 max-w-2xl">
        {/* Logo/Brand */}
        <div className="space-y-4">
          <div className="flex items-center justify-center gap-3">
            <Rocket className="h-12 w-12 text-primary" />
            <h1 className="text-4xl font-bold tracking-tight">DeployFlow</h1>
          </div>
          <p className="text-muted-foreground text-lg">
            Deploy your Python models to production in seconds.
            <br />
            No Kubernetes. No Docker knowledge required.
          </p>
        </div>

        {/* CTA */}
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
          <DialogTrigger asChild>
            <Button size="lg" className="gap-2 text-base px-8 py-6">
              <Plus className="h-5 w-5" />
              Create New Project
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Create New Project</DialogTitle>
              <DialogDescription>
                Give your project a unique name. This will be used as your
                subdomain (e.g., my-model.localhost:8000).
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="name">Project Name</Label>
                <Input
                  id="name"
                  placeholder="my-awesome-model"
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') handleCreateProject();
                  }}
                  className="font-mono"
                />
                {error && (
                  <p className="text-sm text-red-500">{error}</p>
                )}
              </div>
            </div>
            <DialogFooter>
              <Button
                onClick={handleCreateProject}
                disabled={isLoading}
                className="gap-2"
              >
                {isLoading && <Loader2 className="h-4 w-4 animate-spin" />}
                {isLoading ? 'Creating...' : 'Create Project'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Features */}
        <div className="grid grid-cols-3 gap-6 pt-12 text-left">
          <div className="space-y-2">
            <h3 className="font-semibold">⚡ Instant Deploy</h3>
            <p className="text-sm text-muted-foreground">
              Paste your Flask code and deploy in under 30 seconds.
            </p>
          </div>
          <div className="space-y-2">
            <h3 className="font-semibold">🌐 Auto Subdomain</h3>
            <p className="text-sm text-muted-foreground">
              Each project gets its own URL like project.localhost:8000.
            </p>
          </div>
          <div className="space-y-2">
            <h3 className="font-semibold">🐳 Docker Powered</h3>
            <p className="text-sm text-muted-foreground">
              Your code runs in isolated containers automatically.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
