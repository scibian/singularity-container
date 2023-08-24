#
# Copyright (c) 2017-2022, SyLabs, Inc. All rights reserved.
# Copyright (c) 2017, SingularityWare, LLC. All rights reserved.
#
# Copyright (c) 2015-2017, Gregory M. Kurtzer. All rights reserved.
#
# Copyright (c) 2016, The Regents of the University of California, through
# Lawrence Berkeley National Laboratory (subject to receipt of any required
# approvals from the U.S. Dept. of Energy).  All rights reserved.
#
# This software is licensed under a customized 3-clause BSD license.  Please
# consult LICENSE file distributed with the sources of this project regarding
# your rights to use or distribute this software.
#
# NOTICE.  This Software was developed under funding from the U.S. Department of
# Energy and the U.S. Government consequently retains certain rights. As such,
# the U.S. Government has been granted for itself and others acting on its
# behalf a paid-up, nonexclusive, irrevocable, worldwide license in the Software
# to reproduce, distribute copies to the public, prepare derivative works, and
# perform publicly and display publicly, and to permit other to do so.
#
#

# Disable debugsource packages; otherwise it ends up with an empty %files
#   file in debugsourcefiles.list on Fedora
%undefine _debugsource_packages

Name: singularity-ce
Version: 3.11.4
Release: 1%{?dist}
Summary: Application and environment virtualization

# See LICENSE.md for first party code (BSD-3-Clause and LBNL BSD)
# See LICENSE_THIRD_PARTY.md for incorporated code (ASL 2.0)
# See LICENSE_DEPENDENCIES.md for dependencies
# License identifiers taken from: https://fedoraproject.org/wiki/Licensing
License: BSD-3-Clause and LBNL BSD and ASL 2.0

URL: https://www.sylabs.io/singularity/
Source: %{name}-3.11.4.tar.gz

# Note - we do not require Golang. It can be too old in distros, and we assume it
# may be provided outside of a distro package here. This does break building via
# mock. Distro packages derived from this spec need to require it.
BuildRequires: gcc
BuildRequires: make
# Paths to runtime dependencies detected by mconfig, so must be present at build time.
BuildRequires: cryptsetup
%if "%{_target_vendor}" == "suse"
BuildRequires: squashfs
%else
BuildRequires: squashfs-tools
%endif
# Required for building bundled conmon
BuildRequires: libseccomp-devel
BuildRequires: glib2-devel

# crun requirement not satisfied on EL7 or SLES default repos - use runc there.
%if "%{_target_vendor}" == "suse" || 0%{?rhel} < 8
Requires: runc
%else
Requires: crun
%endif
%if "%{_target_vendor}" == "suse"
Requires: squashfs
%else
Requires: shadow-utils
Requires: squashfs-tools
%endif
Requires: cryptsetup

Provides: %{name}-runtime

# Conflicts with non-CE packages
Conflicts: singularity
# Conflicts with Apptainer, which installs the `/usr/bin/singularity` compatibility executable
Conflicts: apptainer
# Conflicts with SingularityPRO basic packaging (not other variants).
Conflicts: singularitypro24
Conflicts: singularitypro25
Conflicts: singularitypro26
Conflicts: singularitypro31
Conflicts: singularitypro35
Conflicts: singularitypro37
Conflicts: singularitypro39

%description
SingularityCE is the Community Edition of Singularity, an open source
container platform designed to be simple, fast, and secure.

%prep
%autosetup -n %{name}-3.11.4

%build
./mconfig -V %{version}-%{release} \
        --prefix=%{_prefix} \
        --exec-prefix=%{_exec_prefix} \
        --bindir=%{_bindir} \
        --sbindir=%{_sbindir} \
        --sysconfdir=%{_sysconfdir} \
        --datadir=%{_datadir} \
        --includedir=%{_includedir} \
        --libdir=%{_libdir} \
        --libexecdir=%{_libexecdir} \
        --localstatedir=%{_localstatedir} \
        --sharedstatedir=%{_sharedstatedir} \
        --mandir=%{_mandir} \
        --infodir=%{_infodir}

%make_build -C builddir old_config= V=

%install
%make_install -C builddir V=

%files
%attr(4755, root, root) %{_libexecdir}/singularity/bin/starter-suid
%{_bindir}/singularity
%{_bindir}/run-singularity
%dir %{_libexecdir}/singularity
%dir %{_libexecdir}/singularity/bin
%{_libexecdir}/singularity/bin/conmon
%{_libexecdir}/singularity/bin/starter
%dir %{_libexecdir}/singularity/cni
%{_libexecdir}/singularity/cni/*
%dir %{_sysconfdir}/singularity
%config(noreplace) %{_sysconfdir}/singularity/*.conf
%config(noreplace) %{_sysconfdir}/singularity/*.toml
%config(noreplace) %{_sysconfdir}/singularity/*.json
%config(noreplace) %{_sysconfdir}/singularity/*.yaml
%config(noreplace) %{_sysconfdir}/singularity/global-pgp-public
%dir %{_sysconfdir}/singularity/cgroups
%config(noreplace) %{_sysconfdir}/singularity/cgroups/*
%dir %{_sysconfdir}/singularity/network
%config(noreplace) %{_sysconfdir}/singularity/network/*
%dir %{_sysconfdir}/singularity/seccomp-profiles
%config(noreplace) %{_sysconfdir}/singularity/seccomp-profiles/*
%dir %{_sysconfdir}/bash_completion.d
%{_sysconfdir}/bash_completion.d/*
%dir %{_localstatedir}/singularity
%dir %{_localstatedir}/singularity/mnt
%dir %{_localstatedir}/singularity/mnt/session
%{_mandir}/man1/singularity*
%license LICENSE.md
%license LICENSE_THIRD_PARTY.md
%license LICENSE_DEPENDENCIES.md
%doc README.md
%doc CHANGELOG.md
%doc CONTRIBUTING.md

%changelog
